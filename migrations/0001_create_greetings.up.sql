CREATE TABLE IF NOT EXISTS greetings (
    id VARCHAR(50) PRIMARY KEY,
    message VARCHAR(255) NOT NULL
);

INSERT INTO greetings (id, message) VALUES ('1', '¡Hola Mundo desde PostgreSQL en GCP!') ON CONFLICT DO NOTHING;
