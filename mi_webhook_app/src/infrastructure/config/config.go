// File: src/infrastructure/config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Mantiene dependencia de godotenv aquí
)

// AppConfig mantiene la configuración de la aplicación.
type AppConfig struct {
	Port                         string
	DiscordWebhookURLDevelopment string
	DiscordWebhookURLTesting     string
	// GithubWebhookSecret string // Descomenta si usas verificación de firma
}

// LoadConfig carga la configuración desde variables de entorno.
func LoadConfig() (*AppConfig, error) {
	// Carga archivo .env primero, ignora error si no se encuentra
	err := godotenv.Load()
	if err != nil {
		// No es un error fatal si .env no existe (podrían estar seteadas en el sistema)
		log.Println("WARNING: Could not load .env file, reading environment variables directly.")
	}

	devURL := os.Getenv("DISCORD_WEBHOOK_URL_DEVELOPMENT")
	if devURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL_DEVELOPMENT environment variable not set")
	}

	testURL := os.Getenv("DISCORD_WEBHOOK_URL_TESTING")
	if testURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL_TESTING environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
		log.Printf("INFO: PORT environment variable not set, using default %s", port)
	}

	// secret := os.Getenv("GITHUB_WEBHOOK_SECRET") // Descomenta si usas verificación
	// if secret == "" {
	//     log.Println("WARNING: GITHUB_WEBHOOK_SECRET environment variable not set. Signature verification disabled.")
	// }

	return &AppConfig{
		Port:                         port,
		DiscordWebhookURLDevelopment: devURL,
		DiscordWebhookURLTesting:     testURL,
		// GithubWebhookSecret: secret, // Descomenta
	}, nil
}