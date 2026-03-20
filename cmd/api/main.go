package main

import (
	"net/http"
	"time"

	"springboard/internal/config"
	"springboard/internal/lib"
	"springboard/internal/middleware"
	"springboard/internal/storage"

	authHandler "springboard/internal/auth/handler"
	authRepo "springboard/internal/auth/repository"
	authService "springboard/internal/auth/service"

	userHandler "springboard/internal/user/handler"
	userRepo "springboard/internal/user/repository"
	userService "springboard/internal/user/service"

	oppHandler "springboard/internal/opportunity/handler"
	oppRepo "springboard/internal/opportunity/repository"
	oppService "springboard/internal/opportunity/service"
)

func main() {
	cfg := config.MustLoadConfig()
	db := storage.MustLoadDatabase(cfg.DatabaseURL)
	defer db.Close()

	jwtManager := lib.NewJWTManager(cfg.SecretKey)

	// USER LAYER
	uRepo := userRepo.NewUserRepository(db)
	uService := userService.NewUserService(uRepo)
	uHandler := userHandler.NewUserHandler(uService)

	// AUTH LAYER
	aRepo := authRepo.NewAuthRepository(db)
	aService := authService.NewAuthService(aRepo, jwtManager)
	aHandler := authHandler.NewAuthHandler(aService)

	// OPPORTUNITY LAYER
	oRepo := oppRepo.NewOpportunityRepository(db)
	oService := oppService.NewOpportunityService(oRepo, uRepo)
	oHandler := oppHandler.NewOpportunityHandler(oService)

	// MIDDLEWARE
	authMW := middleware.CheckTokenMiddleware(jwtManager)

	// ROUTING
	api := http.NewServeMux()

	// PUBLIC ROUTES
	aHandler.RegisterRoutes(api)

	// PRIVATE & PROTECTED ROUTES
	uHandler.RegisterRoutes(api, authMW)
	oHandler.RegisterRoutes(api, authMW)

	// CREATE A GROUP /api
	mainMux := http.NewServeMux()

	mainMux.Handle("/api/", http.StripPrefix("/api", api))

	wrappedMux := middleware.LoggerMiddleware(mainMux)

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      wrappedMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}
