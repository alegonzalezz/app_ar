package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"gcp-serverless-app/auth"
	"gcp-serverless-app/hello"
	"gcp-serverless-app/migrations"
	"gcp-serverless-app/users"

	_ "github.com/lib/pq"
)

func main() {
	dbConnStr := getEnv("DATABASE_URL", "postgres://user_dev:password_dev@localhost:5432/dev_db?sslmode=disable")
	
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

	// Inicializar slice Hello
	helloRepo := hello.NewRepository(db)
	helloUseCase := hello.NewFindGreetingUseCase(helloRepo)
	helloHandler := hello.NewHandler(helloUseCase)
	http.Handle("/hello", helloHandler)

	// Inicializar slice Auth
	authRepo := auth.NewRepository(db)
	createAuthUseCase := auth.NewCreateAuthUseCase(authRepo)
	changePasswordUseCase := auth.NewChangePasswordUseCase(authRepo)
	changePasswordHandler := auth.NewChangePasswordHandler(changePasswordUseCase)
	http.Handle("/auth/change-password", changePasswordHandler)

	// Inicializar slice Users (con dependencia a Auth vía adapter/bridge)
	usersRepo := users.NewRepository(db)
	txManager := users.NewTxManager(db)
	authCreatorBridge := users.NewAuthBridge(createAuthUseCase)
	createUserUseCase := users.NewCreateUserUseCase(txManager, usersRepo, authCreatorBridge)
	usersHandler := users.NewHandler(createUserUseCase)
	http.Handle("/users", usersHandler)

	port := getEnv("PORT", "8080")
	fmt.Printf("Servidor corriendo en el puerto %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
