// File: src/infrastructure/router/router.go
package router

import (
	// --- IMPORTACIONES ACTUALIZADAS (usa tu nombre de módulo) ---
	"mi_webhook_app/src/application"
	"mi_webhook_app/src/infrastructure/handlers"
	// --- FIN IMPORTACIONES ACTUALIZADAS ---

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura el motor Gin.
// Recibe el servicio de aplicación (a través de su puerto) para inyectarlo en el manejador.
func SetupRoutes(engine *gin.Engine, processor application.WebhookProcessor) {

	// Endpoint base para los webhooks entrantes
	webhookGroup := engine.Group("/webhook")
	{
		// Un único endpoint para recibir todos los webhooks de GitHub
		// Pasa el servicio de aplicación (processor) a la fábrica de manejadores.
		// Si necesitaras inyectar config al handler (ej: para firma), lo harías aquí.
		webhookGroup.POST("/github", handlers.GithubWebhookHandler(processor /*, cfg */))
	}

	// Endpoint opcional de health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})
}