package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	appasset "github.com/smetanamolokovich/veylo/internal/application/asset"
	appauth "github.com/smetanamolokovich/veylo/internal/application/auth"
	appfinding "github.com/smetanamolokovich/veylo/internal/application/finding"
	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/bcrypt"
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

	// Wire up dependencies
	// inspections
	inspectionRepo := postgres.NewInspectionRepository(db)
	createInspection := appinspection.NewCreateInspectionUseCase(inspectionRepo)
	getInspection := appinspection.NewGetInspectionUseCase(inspectionRepo)
	listInspections := appinspection.NewListInspectionsUseCase(inspectionRepo)
	transitionInspection := appinspection.NewTransitionInspectionUseCase(inspectionRepo)
	inspectionHandler := handler.NewInspectionHandler(createInspection, listInspections, getInspection, transitionInspection)

	// auth
	userRepo := postgres.NewUserRepository(db)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db)
	registerUC := appauth.NewRegisterUseCase(userRepo, hasher)
	loginUC := appauth.NewLoginUseCase(userRepo, refreshTokenRepo, hasher, jwtManager)
	refreshUC := appauth.NewRefreshTokenUseCase(refreshTokenRepo, userRepo, jwtManager, hasher)
	authHandler := handler.NewAuthHandler(registerUC, loginUC, refreshUC)

	// assets
	assetRepo := postgres.NewAssetRepository(db)
	createVehicleUC := appasset.NewCreateVehicleAssetUseCase(assetRepo)
	getAssetUC := appasset.NewGetAssetUseCase(assetRepo)
	assetHandler := handler.NewAssetHandler(createVehicleUC, getAssetUC)

	// findings
	findingRepo := postgres.NewFindingRepository(db)
	createFindingUC := appfinding.NewCreateFindingUseCase(findingRepo)
	listFindingsUC := appfinding.NewListFindingsUseCase(findingRepo)
	assessFindingUC := appfinding.NewAssessFindingUseCase(findingRepo)
	findingHandler := handler.NewFindingHandler(createFindingUC, listFindingsUC, assessFindingUC)

	router := httpinterface.NewRouter(inspectionHandler, authHandler, assetHandler, findingHandler, jwtManager)

	addr := ":8080"
	log.Info("starting server", "addr", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error("server error", "err", err)
		os.Exit(1)
	}
}
