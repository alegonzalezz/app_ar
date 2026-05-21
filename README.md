# app_ar - Go Serverless Architecture

Este proyecto está diseñado bajo los patrones de **Vertical Slice Architecture** y **Clean Architecture** enfocado a despliegues Serverless eficientes (GCP Cloud Run / Functions) utilizando código nativo en Go.

## 🚀 Cómo inicializar el entorno local

1. **Levantar la base de datos PostgreSQL:**
   docker-compose up -d

2. **Instalar dependencias de Go:**
   go mod tidy

3. **Correr la aplicación de forma nativa:**
   go run main.go

4. **Probar el endpoint:**
   curl "http://localhost:8080/hello?id=1"

## 🛠️ Especificación para IAs (Cursor / Claude / OpenCode)

El archivo PROMPT_AI_SPEC.md contiene las reglas estrictas de diseño de este repositorio. Siempre que quieras agregar un endpoint, use case, o cliente externo, pasale ese archivo a la IA elegida para garantizar el cumplimiento del patrón arquitectónico.
