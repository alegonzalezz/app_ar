package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"gcp-serverless-app/hello"

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

	repo := hello.NewRepository(db)
	useCase := hello.NewFindGreetingUseCase(repo)
	helloHandler := hello.NewHandler(useCase)

	http.Handle("/hello", helloHandler)

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
