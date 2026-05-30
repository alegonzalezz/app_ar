package config

import "os"

// GetEnv retorna el valor de la variable de entorno o el fallback.
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
