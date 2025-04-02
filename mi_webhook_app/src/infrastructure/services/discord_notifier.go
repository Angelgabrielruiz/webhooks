// File: src/infrastructure/services/discord_notifier.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// --- IMPORTACIONES ACTUALIZADAS (usa tu nombre de módulo) ---
	"mi_webhook_app/src/application"
	"mi_webhook_app/src/infrastructure/config"
	// --- FIN IMPORTACIONES ACTUALIZADAS ---
)

// discordNotifier es la implementación concreta para enviar notificaciones a Discord.
type discordNotifier struct {
	config *config.AppConfig
}

// NewDiscordNotifier crea un nuevo adaptador implementando application.NotificationService.
func NewDiscordNotifier(cfg *config.AppConfig) application.NotificationService {
	return &discordNotifier{
		config: cfg,
	}
}

// SendNotification implementa la interfaz application.NotificationService.
func (n *discordNotifier) SendNotification(channelType string, payload application.DiscordPayload) error {
	webhookURL := n.getWebhookURL(channelType)
	if webhookURL == "" {
		// Es importante loguear pero también retornar error para que la app sepa que falló
		log.Printf("ERROR: No Discord webhook URL configured for channel type: %s", channelType)
		return fmt.Errorf("no webhook URL configured for channel type '%s'", channelType)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: Marshalling Discord payload: %v", err)
		return fmt.Errorf("error marshalling discord payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("ERROR: Sending message to Discord (%s): %v", webhookURL, err)
		return fmt.Errorf("error sending http request to discord url %s: %w", webhookURL, err)
	}
	defer resp.Body.Close() // Siempre cierra el cuerpo

	// Verifica el código de estado de Discord
	if resp.StatusCode >= 300 {
		// Intenta leer el cuerpo de la respuesta de Discord para más detalles
		bodyBytes := new(bytes.Buffer)
		_, readErr := bodyBytes.ReadFrom(resp.Body)
		if readErr != nil {
			log.Printf("ERROR: Reading Discord error response body: %v", readErr)
		}
		log.Printf("ERROR: Discord webhook (%s) returned non-success status: %s. Body: %s", webhookURL, resp.Status, bodyBytes.String())
		// Retorna un error que indica el fallo
		return fmt.Errorf("discord webhook (%s) failed with status %s", webhookURL, resp.Status)
	}

	log.Printf("INFO: Successfully sent notification to Discord channel type '%s'", channelType)
	return nil // Éxito
}

// getWebhookURL recupera la URL apropiada basada en el tipo de canal lógico.
func (n *discordNotifier) getWebhookURL(channelType string) string {
	switch channelType {
	case "development":
		return n.config.DiscordWebhookURLDevelopment
	case "testing":
		return n.config.DiscordWebhookURLTesting
	default:
		log.Printf("WARNING: Unknown channel type requested: %s", channelType)
		return "" // Retorna vacío si el tipo no es conocido
	}
}