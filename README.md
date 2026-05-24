# app_ar - Go Serverless Architecture

Este proyecto está diseñado bajo los patrones de **Vertical Slice Architecture** y **Clean Architecture** enfocado a despliegues Serverless eficientes (GCP Cloud Run / Functions) utilizando código nativo en Go.

---

## 🚀 Cómo inicializar el entorno local

1. **Levantar la base de datos PostgreSQL:**
   ```bash
   docker-compose up -d postgres
   ```
   *Nota: No es necesario levantar el servicio `app` en docker si vas a correr el código de forma nativa.*

2. **Instalar dependencias de Go:**
   ```bash
   go mod tidy
   ```

3. **Correr la aplicación de forma nativa:**
   ```bash
   go run main.go
   ```
   *Al iniciar la aplicación por primera vez, el sistema de migraciones nativo detectará la base de datos y ejecutará automáticamente todos los scripts pendientes dentro del directorio `migrations/`.*

4. **Probar los endpoints:**
   - **Obtener saludo (Hello Slice):**
     ```bash
     curl "http://localhost:8080/hello?id=1"
     ```
   - **Crear Usuario (Users Slice - Transaccional con Auth):**
     ```bash
     curl -i -X POST http://localhost:8080/users \
       -H "Content-Type: application/json" \
       -d '{"name": "Juan Perez", "email": "juan.perez@example.com", "password": "SuperSecretPassword123"}'
     ```

---

## 🔄 Sistema de Migraciones e Historial

El proyecto cuenta con un sistema de migraciones autogestionado en Go. Las migraciones se guardan como archivos SQL secuenciales dentro de la carpeta `migrations/`.

Para auditar y llevar el control, el sistema mantiene la tabla **`schema_migrations`** en la base de datos, la cual almacena:
- `version` (TEXT/VARCHAR): Nombre del archivo de migración ejecutado.
- `applied_at` (TIMESTAMP WITH TIME ZONE): Fecha y hora exacta de cuándo fue aplicada la migración.

Puedes consultar el historial ejecutando:
```bash
docker exec -it local_postgres psql -U user_dev -d dev_db -c "SELECT * FROM schema_migrations;"
```

---

## 🧹 Cómo reiniciar la base de datos desde cero

Si necesitas realizar pruebas limpias y ejecutar todas las migraciones nuevamente:

1. **Detener contenedores y eliminar volúmenes asociados:**
   ```bash
   docker-compose down -v
   ```
   *El flag `-v` es crítico aquí, ya que destruye el volumen de datos persistente de Postgres (`pgdata`).*

2. **Volver a levantar la base de datos limpia:**
   ```bash
   docker-compose up -d postgres
   ```

3. **Iniciar el servidor nativo para que corra las migraciones automáticamente:**
   ```bash
   go run main.go
   ```

---

## 🛠️ Especificación para IAs (Cursor / Claude / OpenCode)

El archivo PROMPT_AI_SPEC.md contiene las reglas estrictas de diseño de este repositorio. Siempre que quieras agregar un endpoint, use case, o cliente externo, pasale ese archivo a la IA elegida para garantizar el cumplimiento del patrón arquitectónico.

