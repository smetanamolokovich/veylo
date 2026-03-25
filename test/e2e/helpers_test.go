package e2e_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	appasset "github.com/smetanamolokovich/veylo/internal/application/asset"
	appauth "github.com/smetanamolokovich/veylo/internal/application/auth"
	appfinding "github.com/smetanamolokovich/veylo/internal/application/finding"
	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	apporg "github.com/smetanamolokovich/veylo/internal/application/organization"
	appworkflow "github.com/smetanamolokovich/veylo/internal/application/workflow"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/bcrypt"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/postgres"
	httpinterface "github.com/smetanamolokovich/veylo/internal/interface/http"
	"github.com/smetanamolokovich/veylo/internal/interface/http/handler"
	"github.com/smetanamolokovich/veylo/pkg/jwt"
)

// testServer holds a running test HTTP server and its dependencies.
type testServer struct {
	server *httptest.Server
	token  string // access token after signup
}

func newTestServer(t *testing.T) (*testServer, func()) {
	t.Helper()
	ctx := context.Background()

	// Start postgres container
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("veylo_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		testcontainers.WithAdditionalWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	require.NoError(t, db.Ping())

	// Run migrations
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(thisFile), "..", "..", "migrations")

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	require.NoError(t, err)
	require.NoError(t, m.Up())

	// Wire dependencies
	jwtManager := jwt.NewManager("test-secret")
	hasher := bcrypt.NewPasswordHasher()

	workflowRepo := postgres.NewWorkflowRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	assetRepo := postgres.NewAssetRepository(db)
	findingRepo := postgres.NewFindingRepository(db)
	inspectionRepo := postgres.NewInspectionRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	// Auth
	registerUC := appauth.NewRegisterUseCase(userRepo, refreshTokenRepo, hasher, jwtManager)
	loginUC := appauth.NewLoginUseCase(userRepo, refreshTokenRepo, hasher, jwtManager)
	refreshUC := appauth.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, jwtManager, hasher)
	signupUC := appauth.NewSignupUseCase(orgRepo, workflowRepo, userRepo, refreshTokenRepo, hasher, jwtManager)
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC, signupUC)

	// Workflow
	createWorkflowUC := appworkflow.NewCreateWorkflowUseCase(workflowRepo)
	getWorkflowUC := appworkflow.NewGetWorkflowUseCase(workflowRepo)
	addStatusUC := appworkflow.NewAddStatusUseCase(workflowRepo)
	addTransitionUC := appworkflow.NewAddTransitionUseCase(workflowRepo)
	workflowHandler := handler.NewWorkflowHandler(createWorkflowUC, getWorkflowUC, addStatusUC, addTransitionUC)

	// Org
	createOrgUC := apporg.NewCreateOrganizationUseCase(orgRepo, workflowRepo, userRepo, jwtManager)
	completeOnboardingUC := apporg.NewCompleteOnboardingUseCase(orgRepo)
	orgHandler := handler.NewOrganizationHandler(orgRepo, createOrgUC, completeOnboardingUC)

	// Assets
	createVehicleUC := appasset.NewCreateVehicleAssetUseCase(assetRepo)
	getAssetUC := appasset.NewGetAssetUseCase(assetRepo)
	assetHandler := handler.NewAssetHandler(createVehicleUC, getAssetUC)

	// Findings
	createFindingUC := appfinding.NewCreateFindingUseCase(findingRepo)
	listFindingsUC := appfinding.NewListFindingsUseCase(findingRepo)
	assessFindingUC := appfinding.NewAssessFindingUseCase(findingRepo)
	findingHandler := handler.NewFindingHandler(createFindingUC, listFindingsUC, assessFindingUC)

	// Inspections (no S3 in tests)
	createInspection := appinspection.NewCreateInspectionUseCase(inspectionRepo, workflowRepo)
	getInspection := appinspection.NewGetInspectionUseCase(inspectionRepo)
	listInspections := appinspection.NewListInspectionsUseCase(inspectionRepo)
	transitionInspection := appinspection.NewTransitionInspectionUseCase(inspectionRepo, workflowRepo, nil)
	inspectionHandler := handler.NewInspectionHandler(createInspection, listInspections, getInspection, transitionInspection, reportRepo)

	router := httpinterface.NewRouter(inspectionHandler, authHandler, assetHandler, findingHandler, workflowHandler, orgHandler, jwtManager)
	srv := httptest.NewServer(router)

	cleanup := func() {
		srv.Close()
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return &testServer{server: srv}, cleanup
}

// do sends an HTTP request and decodes the JSON response into out (if not nil).
func (ts *testServer) do(t *testing.T, method, path string, body any, out any) *http.Response {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	req, err := http.NewRequest(method, ts.server.URL+path, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if ts.token != "" {
		req.Header.Set("Authorization", "Bearer "+ts.token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	if out != nil {
		require.NoError(t, json.NewDecoder(resp.Body).Decode(out))
		resp.Body.Close()
	}

	return resp
}

func (ts *testServer) mustJSON(t *testing.T, resp *http.Response, expectedStatus int) map[string]any {
	t.Helper()
	require.Equal(t, expectedStatus, resp.StatusCode, "unexpected status for %s", resp.Request.URL)
	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	resp.Body.Close()
	return result
}

func ptr[T any](v T) *T { return &v }

func jsonBody(pairs ...any) map[string]any {
	m := make(map[string]any)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[fmt.Sprintf("%v", pairs[i])] = pairs[i+1]
	}
	return m
}
