package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	appasset "github.com/smetanamolokovich/veylo/internal/application/asset"
	appauth "github.com/smetanamolokovich/veylo/internal/application/auth"
	appfinding "github.com/smetanamolokovich/veylo/internal/application/finding"
	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	apporg "github.com/smetanamolokovich/veylo/internal/application/organization"
	appreport "github.com/smetanamolokovich/veylo/internal/application/report"
	appworkflow "github.com/smetanamolokovich/veylo/internal/application/workflow"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/bcrypt"
	infrapdf "github.com/smetanamolokovich/veylo/internal/infrastructure/pdf"
	infraS3 "github.com/smetanamolokovich/veylo/internal/infrastructure/s3"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/postgres"
	httpinterface "github.com/smetanamolokovich/veylo/internal/interface/http"
	"github.com/smetanamolokovich/veylo/internal/interface/http/handler"
	"github.com/smetanamolokovich/veylo/pkg/jwt"
	"github.com/smetanamolokovich/veylo/pkg/logger"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	log := logger.New(env)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:6543/veylo?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("database is not reachable", "err", err)
		os.Exit(1)
	}
	log.Info("database connected")

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		if env == "production" {
			log.Error("JWT_SECRET environment variable is required in production")
			os.Exit(1)
		} else {
			secret = "dev-secret"
			log.Warn("using default JWT secret in non-production environment")
		}
	}

	jwtManager := jwt.NewManager(secret)
	hasher := bcrypt.NewPasswordHasher()

	// Repositories
	workflowRepo := postgres.NewWorkflowRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	assetRepo := postgres.NewAssetRepository(db)
	findingRepo := postgres.NewFindingRepository(db)
	inspectionRepo := postgres.NewInspectionRepository(db)
	reportRepo := postgres.NewReportRepository(db)

	// Report generation (optional — requires S3_BUCKET env var)
	var generateReportUC *appreport.GenerateReportUseCase
	s3Bucket := os.Getenv("S3_BUCKET")
	s3BaseURL := os.Getenv("S3_BASE_URL")
	if s3Bucket != "" {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Error("failed to load AWS config", "err", err)
			os.Exit(1)
		}
		s3Client := awss3.NewFromConfig(cfg)
		uploader := infraS3.NewUploader(s3Client, s3Bucket, s3BaseURL)
		pdfGenerator := infrapdf.NewVehicleReportGenerator()
		generateReportUC = appreport.NewGenerateReportUseCase(
			inspectionRepo, assetRepo, findingRepo, orgRepo, reportRepo, pdfGenerator, uploader,
		)
		log.Info("S3 report generation enabled", "bucket", s3Bucket)
	} else {
		log.Warn("S3_BUCKET not set — PDF report generation disabled")
	}

	// Workflow
	createWorkflowUC := appworkflow.NewCreateWorkflowUseCase(workflowRepo)
	getWorkflowUC := appworkflow.NewGetWorkflowUseCase(workflowRepo)
	addStatusUC := appworkflow.NewAddStatusUseCase(workflowRepo)
	addTransitionUC := appworkflow.NewAddTransitionUseCase(workflowRepo)
	workflowHandler := handler.NewWorkflowHandler(createWorkflowUC, getWorkflowUC, addStatusUC, addTransitionUC)

	// Organizations
	createOrgUC := apporg.NewCreateOrganizationUseCase(orgRepo, workflowRepo, userRepo, jwtManager)
	completeOnboardingUC := apporg.NewCompleteOnboardingUseCase(orgRepo)
	orgHandler := handler.NewOrganizationHandler(orgRepo, createOrgUC, completeOnboardingUC)

	// Auth
	registerUC := appauth.NewRegisterUseCase(userRepo, refreshTokenRepo, hasher, jwtManager)
	loginUC := appauth.NewLoginUseCase(userRepo, refreshTokenRepo, hasher, jwtManager)
	refreshUC := appauth.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, jwtManager, hasher)
	signupUC := appauth.NewSignupUseCase(orgRepo, workflowRepo, userRepo, refreshTokenRepo, hasher, jwtManager)
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC, signupUC)

	// Assets
	createVehicleUC := appasset.NewCreateVehicleAssetUseCase(assetRepo)
	getAssetUC := appasset.NewGetAssetUseCase(assetRepo)
	assetHandler := handler.NewAssetHandler(createVehicleUC, getAssetUC)

	// Findings
	createFindingUC := appfinding.NewCreateFindingUseCase(findingRepo)
	listFindingsUC := appfinding.NewListFindingsUseCase(findingRepo)
	assessFindingUC := appfinding.NewAssessFindingUseCase(findingRepo)
	findingHandler := handler.NewFindingHandler(createFindingUC, listFindingsUC, assessFindingUC)

	// Inspections
	createInspection := appinspection.NewCreateInspectionUseCase(inspectionRepo, workflowRepo)
	getInspection := appinspection.NewGetInspectionUseCase(inspectionRepo)
	listInspections := appinspection.NewListInspectionsUseCase(inspectionRepo)
	transitionInspection := appinspection.NewTransitionInspectionUseCase(inspectionRepo, workflowRepo, generateReportUC)
	inspectionHandler := handler.NewInspectionHandler(createInspection, listInspections, getInspection, transitionInspection, reportRepo)

	router := httpinterface.NewRouter(inspectionHandler, authHandler, assetHandler, findingHandler, workflowHandler, orgHandler, jwtManager)

	addr := ":8080"
	log.Info("starting server", "addr", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error("server error", "err", err)
		os.Exit(1)
	}
}
