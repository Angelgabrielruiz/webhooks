// File: main.go
package main

import (
	"log"

	// --- IMPORTACIONES ACTUALIZADAS (usa tu nombre de módulo) ---
	"mi_webhook_app/src/application"
	"mi_webhook_app/src/infrastructure/config"
	"mi_webhook_app/src/infrastructure/router"
	"mi_webhook_app/src/infrastructure/services"
	// --- FIN IMPORTACIONES ACTUALIZADAS ---

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration (Infrastructure)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ERROR: Failed to load configuration: %v", err)
	}

	// 2. Initialize Driven Adapters (Infrastructure)
	// Crea el adaptador concreto del notificador Discord
	discordNotifier := services.NewDiscordNotifier(cfg)

	// 3. Initialize Application Service (Core)
	// Crea el servicio de aplicación central, inyectando el adaptador notificador
	// a través del puerto de interfaz application.NotificationService.
	webhookService := application.NewWebhookService(discordNotifier)

	// 4. Initialize Driving Adapters (Infrastructure)
	gin.SetMode(gin.ReleaseMode) // O gin.DebugMode
	engine := gin.Default()
	// Configura rutas, inyectando el servicio de aplicación (webhookService)
	// que cumple con el puerto application.WebhookProcessor.
	router.SetupRoutes(engine, webhookService)

	// 5. Start the Server (Infrastructure)
	log.Printf("INFO: Server starting on port %s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("FATAL: Failed to run server: %v", err)
	}
}