package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"gcp-serverless-app/internal/config"
	"gcp-serverless-app/migrations"

	// Infrastructure
	authInfra "gcp-serverless-app/internal/auth/infrastructure"
	greetInfra "gcp-serverless-app/internal/greeting/infrastructure"
	pg "gcp-serverless-app/internal/shared/platform/postgres"
	userInfra "gcp-serverless-app/internal/user/infrastructure"

	// Application
	authApp "gcp-serverless-app/internal/auth/application"
	greetApp "gcp-serverless-app/internal/greeting/application"
	userApp "gcp-serverless-app/internal/user/application"

	// Bridges
	"gcp-serverless-app/internal/shared/bridge"

	// HTTP Handlers
	"gcp-serverless-app/internal/shared/platform/http/handlers"
)

func main() {
	dbConnStr := config.GetEnv("DATABASE_URL", "postgres://user_dev:password_dev@localhost:5432/dev_db?sslmode=disable")

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Error abriendo conexión a la DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("Advertencia: No se pudo conectar a la DB: %v", err)
	}

	// Ejecutar migraciones en el arranque
	ctx := context.Background()
	if err := migrations.Run(ctx, db); err != nil {
		log.Fatalf("Error ejecutando migraciones: %v", err)
	}

	// === INFRAESTRUCTURA ===
	txManager := pg.NewTxManager(db)
	hasher := authInfra.NewSHA256Hasher()

	authRepo := authInfra.NewPostgresRepository(db)
	userRepo := userInfra.NewPostgresRepository(db)
	greetRepo := greetInfra.NewPostgresRepository(db)

	// === CASOS DE USO ===
	createAuthUC := authApp.NewCreateAuthUseCase(authRepo, hasher)
	changePassUC := authApp.NewChangePasswordUseCase(authRepo, hasher)

	// Bridges inter-módulo
	authCreator := bridge.NewAuthCreatorBridge(createAuthUC)
	userProvider := bridge.NewUserProviderBridge(userRepo)

	createUserUC := userApp.NewCreateUserUseCase(txManager, userRepo, authCreator)
	loginUC := authApp.NewLoginUseCase(authRepo, userProvider, hasher)
	findGreetingUC := greetApp.NewFindGreetingUseCase(greetRepo)

	// === HTTP HANDLERS ===
	http.Handle("/users", handlers.NewUserHandler(createUserUC))
	http.Handle("/auth/login", handlers.NewLoginHandler(loginUC))
	http.Handle("/auth/change-password", handlers.NewChangePasswordHandler(changePassUC))
	http.Handle("/hello", handlers.NewGreetingHandler(findGreetingUC))

	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Servidor corriendo en el puerto %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
