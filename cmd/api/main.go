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
	adminInfra "gcp-serverless-app/internal/administrative/infrastructure"
	appointmentInfra "gcp-serverless-app/internal/appointment/infrastructure"
	authInfra "gcp-serverless-app/internal/auth/infrastructure"
	customerInfra "gcp-serverless-app/internal/customer/infrastructure"
	greetInfra "gcp-serverless-app/internal/greeting/infrastructure"
	pg "gcp-serverless-app/internal/shared/platform/postgres"
	taskInfra "gcp-serverless-app/internal/task/infrastructure"
	userInfra "gcp-serverless-app/internal/user/infrastructure"
	visitInfra "gcp-serverless-app/internal/visit/infrastructure"
	workerInfra "gcp-serverless-app/internal/worker/infrastructure"

	// Application
	adminApp "gcp-serverless-app/internal/administrative/application"
	appointmentApp "gcp-serverless-app/internal/appointment/application"
	authApp "gcp-serverless-app/internal/auth/application"
	customerApp "gcp-serverless-app/internal/customer/application"
	greetApp "gcp-serverless-app/internal/greeting/application"
	taskApp "gcp-serverless-app/internal/task/application"
	userApp "gcp-serverless-app/internal/user/application"
	visitApp "gcp-serverless-app/internal/visit/application"
	workerApp "gcp-serverless-app/internal/worker/application"

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
	customerRepo := customerInfra.NewPostgresRepository(db)
	workerRepo := workerInfra.NewPostgresRepository(db)
	taskRepo := taskInfra.NewPostgresRepository(db)
	appointmentRepo := appointmentInfra.NewPostgresRepository(db)
	adminRepo := adminInfra.NewPostgresRepository(db)
	visitRepo := visitInfra.NewPostgresRepository(db)

	// === CASOS DE USO ===
	createAuthUC := authApp.NewCreateAuthUseCase(authRepo, hasher)
	changePassUC := authApp.NewChangePasswordUseCase(authRepo, hasher)

	// Bridges inter-módulo
	authCreator := bridge.NewAuthCreatorBridge(createAuthUC)
	userProvider := bridge.NewUserProviderBridge(userRepo)

	createUserUC := userApp.NewCreateUserUseCase(txManager, userRepo, authCreator)
	loginUC := authApp.NewLoginUseCase(authRepo, userProvider, hasher)
	findGreetingUC := greetApp.NewFindGreetingUseCase(greetRepo)

	createCustomerUC := customerApp.NewCreateCustomerUseCase(customerRepo)
	getCustomerUC := customerApp.NewGetCustomerUseCase(customerRepo)
	listCustomersUC := customerApp.NewListCustomersUseCase(customerRepo)
	updateCustomerUC := customerApp.NewUpdateCustomerUseCase(customerRepo)
	deleteCustomerUC := customerApp.NewDeleteCustomerUseCase(customerRepo)

	createWorkerUC := workerApp.NewCreateWorkerUseCase(workerRepo)
	getWorkerUC := workerApp.NewGetWorkerUseCase(workerRepo)
	listWorkersUC := workerApp.NewListWorkersUseCase(workerRepo)
	updateWorkerUC := workerApp.NewUpdateWorkerUseCase(workerRepo)
	deleteWorkerUC := workerApp.NewDeleteWorkerUseCase(workerRepo)

	createTaskUC := taskApp.NewCreateTaskUseCase(taskRepo)
	getTaskUC := taskApp.NewGetTaskUseCase(taskRepo)
	listTasksUC := taskApp.NewListTasksUseCase(taskRepo)
	updateTaskUC := taskApp.NewUpdateTaskUseCase(taskRepo)
	deleteTaskUC := taskApp.NewDeleteTaskUseCase(taskRepo)
	updateTaskStatusUC := taskApp.NewUpdateTaskStatusUseCase(taskRepo)

	createAppointmentUC := appointmentApp.NewCreateAppointmentUseCase(appointmentRepo)
	getAppointmentUC := appointmentApp.NewGetAppointmentUseCase(appointmentRepo)
	listAppointmentsUC := appointmentApp.NewListAppointmentsUseCase(appointmentRepo)
	updateAppointmentUC := appointmentApp.NewUpdateAppointmentUseCase(appointmentRepo)
	deleteAppointmentUC := appointmentApp.NewDeleteAppointmentUseCase(appointmentRepo)
	updateAppointmentStatusUC := appointmentApp.NewUpdateAppointmentStatusUseCase(appointmentRepo)
	assignTaskUC := appointmentApp.NewAssignTaskUseCase(appointmentRepo)
	unassignTaskUC := appointmentApp.NewUnassignTaskUseCase(appointmentRepo)
	getTasksByAppointmentUC := appointmentApp.NewGetTasksByAppointmentUseCase(appointmentRepo)
	getAppointmentsByTaskUC := appointmentApp.NewGetAppointmentsByTaskUseCase(appointmentRepo)

	createAdminUC := adminApp.NewCreateAdministrativeUseCase(adminRepo)
	getAdminUC := adminApp.NewGetAdministrativeUseCase(adminRepo)
	listAdminsUC := adminApp.NewListAdministrativesUseCase(adminRepo)
	updateAdminUC := adminApp.NewUpdateAdministrativeUseCase(adminRepo)
	deleteAdminUC := adminApp.NewDeleteAdministrativeUseCase(adminRepo)

	createVisitUC := visitApp.NewCreateVisitUseCase(visitRepo)
	getVisitUC := visitApp.NewGetVisitUseCase(visitRepo)
	listVisitsUC := visitApp.NewListVisitsUseCase(visitRepo)
	updateVisitUC := visitApp.NewUpdateVisitUseCase(visitRepo)
	deleteVisitUC := visitApp.NewDeleteVisitUseCase(visitRepo)
	updateVisitStatusUC := visitApp.NewUpdateVisitStatusUseCase(visitRepo)
	assignVisitTaskUC := visitApp.NewAssignTaskUseCase(visitRepo)
	unassignVisitTaskUC := visitApp.NewUnassignTaskUseCase(visitRepo)
	getVisitTasksUC := visitApp.NewGetVisitTasksUseCase(visitRepo)

	// === HTTP HANDLERS ===
	http.Handle("POST /users", handlers.NewUserHandler(createUserUC))
	http.Handle("POST /auth/login", handlers.NewLoginHandler(loginUC))
	http.Handle("POST /auth/change-password", handlers.NewChangePasswordHandler(changePassUC))
	http.Handle("GET /hello", handlers.NewGreetingHandler(findGreetingUC))

	customerHandler := handlers.NewCustomerHandler(createCustomerUC, getCustomerUC, listCustomersUC, updateCustomerUC, deleteCustomerUC)
	customerHandler.RegisterRoutes(http.DefaultServeMux)

	adminHandler := handlers.NewAdministrativeHandler(createAdminUC, getAdminUC, listAdminsUC, updateAdminUC, deleteAdminUC)
	adminHandler.RegisterRoutes(http.DefaultServeMux)

	workerHandler := handlers.NewWorkerHandler(createWorkerUC, getWorkerUC, listWorkersUC, updateWorkerUC, deleteWorkerUC)
	workerHandler.RegisterRoutes(http.DefaultServeMux)

	taskHandler := handlers.NewTaskHandler(createTaskUC, getTaskUC, listTasksUC, updateTaskUC, deleteTaskUC, updateTaskStatusUC)
	taskHandler.RegisterRoutes(http.DefaultServeMux)

	appointmentHandler := handlers.NewAppointmentHandler(
		createAppointmentUC, getAppointmentUC, listAppointmentsUC, updateAppointmentUC, deleteAppointmentUC,
		updateAppointmentStatusUC, assignTaskUC, unassignTaskUC, getTasksByAppointmentUC, getAppointmentsByTaskUC,
	)
	appointmentHandler.RegisterRoutes(http.DefaultServeMux)

	visitHandler := handlers.NewVisitHandler(
		createVisitUC, getVisitUC, listVisitsUC, updateVisitUC, deleteVisitUC,
		updateVisitStatusUC, assignVisitTaskUC, unassignVisitTaskUC, getVisitTasksUC,
	)
	visitHandler.RegisterRoutes(http.DefaultServeMux)

	customerHistoryHandler := handlers.NewCustomerHistoryHandler(listVisitsUC)
	customerHistoryHandler.RegisterRoutes(http.DefaultServeMux)

	workerVisitsHandler := handlers.NewWorkerVisitsHandler(listVisitsUC)
	workerVisitsHandler.RegisterRoutes(http.DefaultServeMux)

	port := config.GetEnv("PORT", "8080")
	fmt.Printf("Servidor corriendo en el puerto %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
