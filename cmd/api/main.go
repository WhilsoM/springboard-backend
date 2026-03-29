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

	appHandler "springboard/internal/application/handler"
	appRepo "springboard/internal/application/repository"
	appService "springboard/internal/application/service"

	adminHandler "springboard/internal/admin/handler"
	adminRepo "springboard/internal/admin/repository"
	adminService "springboard/internal/admin/service"

	"github.com/rs/cors"
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

	// APPLICATION LAYER
	appRepo := appRepo.NewApplicationRepository(db)
	appService := appService.NewApplicationService(appRepo)
	appHandler := appHandler.NewApplicationHandler(appService)

	// ADMIN LAYER
	adminRepo := adminRepo.NewAdminRepository(db)
	adminService := adminService.NewAdminService(adminRepo)
	adminHandler := adminHandler.NewAdminHandler(adminService)

	// MIDDLEWARE
	authMW := middleware.CheckTokenMiddleware(jwtManager)

	// ROUTING
	api := http.NewServeMux()

	// PUBLIC ROUTES
	aHandler.RegisterRoutes(api)

	// PRIVATE & PROTECTED ROUTES
	uHandler.RegisterRoutes(api, authMW)
	oHandler.RegisterRoutes(api, authMW)
	appHandler.RegisterRoutes(api, authMW)
	adminHandler.RegisterRoutes(api, authMW)

	// CREATE A GROUP /api
	mainMux := http.NewServeMux()

	mainMux.Handle("/api/", http.StripPrefix("/api", api))

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://your-production-app.ru", "http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            true,
	}
	c := cors.New(corsOptions)

	wrappedMux := middleware.LoggerMiddleware(mainMux)
	finalHandler := c.Handler(wrappedMux)

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      finalHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		panic("failed to listen and serve: " + err.Error())
	}
}
