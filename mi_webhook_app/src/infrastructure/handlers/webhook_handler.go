// File: src/infrastructure/handlers/webhook_handler.go
package handlers

import (
	"fmt"
	"log"
	"net/http"

	// --- IMPORTACIÓN ACTUALIZADA (usa tu nombre de módulo) ---
	"mi_webhook_app/src/application"
	// --- FIN IMPORTACIÓN ACTUALIZADA ---

	"github.com/gin-gonic/gin"
	// Descomenta para verificación de firma
	// "crypto/hmac"
	// "crypto/sha256"
	// "encoding/hex"
	// "strings"
	// "mi_webhook_app/src/infrastructure/config" // Necesitarías config si verificas firma
)

// GithubWebhookHandler crea una función manejadora de Gin.
// Depende del servicio de aplicación (procesador de casos de uso) a través de su interfaz de puerto.
// Si necesitas verificación de firma, también necesitarías inyectar *config.AppConfig aquí.
func GithubWebhookHandler(processor application.WebhookProcessor /* , cfg *config.AppConfig */) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Headers estándar de GitHub
		eventType := ctx.GetHeader("X-GitHub-Event")
		deliveryID := ctx.GetHeader("X-GitHub-Delivery")
		// signature := ctx.GetHeader("X-Hub-Signature-256") // Para verificación

		log.Printf("INFO: Webhook received: Event=%s, DeliveryID=%s", eventType, deliveryID)

		// Leer payload crudo
		payload, err := ctx.GetRawData()
		if err != nil {
			log.Printf("ERROR: Reading request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Error reading request body"})
			return
		}

		// --- Verificación de Firma (Opcional pero Recomendado) ---
		/*
		   if cfg.GithubWebhookSecret != "" {
		       if !isValidSignature(signature, cfg.GithubWebhookSecret, payload) {
		           log.Printf("WARNING: Invalid webhook signature. DeliveryID: %s", deliveryID)
		           ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid signature"})
		           return
		       }
		       log.Println("INFO: Webhook signature verified.")
		   } else {
		       log.Println("WARNING: Webhook signature verification skipped (GITHUB_WEBHOOK_SECRET not set).")
		   }
		*/

		var processingErr error // Para capturar error de la capa de aplicación

		// Dirige al método apropiado del servicio de aplicación basado en el evento
		switch eventType {
		case "pull_request":
			log.Printf("INFO: Processing 'pull_request' event...")
			// Llama al servicio de aplicación a través del puerto de interfaz
			processingErr = processor.ProcessPullRequestEvent(payload)
		case "workflow_run":
			log.Printf("INFO: Processing 'workflow_run' event...")
			// Llama al servicio de aplicación a través del puerto de interfaz
			processingErr = processor.ProcessWorkflowRunEvent(payload)
		// Añade más casos aquí si manejas otros eventos (ej: push, issues)
		default:
			// Evento recibido pero no manejado por esta aplicación
			log.Printf("INFO: Ignoring unhandled event type: %s", eventType)
			ctx.JSON(http.StatusOK, gin.H{"status": "received", "message": "Event received but type is not handled"})
			return // Importante retornar aquí para no seguir a la lógica de error/éxito
		}

		// Responde a GitHub basado en el resultado de la lógica de aplicación
		if processingErr != nil {
			// Loguea el error específico de la capa de aplicación
			log.Printf("ERROR: Processing event '%s': %v", eventType, processingErr)
			// Retorna un error genérico del servidor al cliente (GitHub)
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": fmt.Sprintf("Error processing event '%s'", eventType)})
		} else {
			// Éxito en el procesamiento (incluso si no se envió notificación por lógica interna)
			log.Printf("INFO: Event '%s' processed successfully. DeliveryID: %s", eventType, deliveryID)
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": fmt.Sprintf("Event '%s' processed successfully", eventType)})
		}
	}
}

/* // Descomenta y ajusta si usas verificación de firma
func isValidSignature(ghSignature, secret string, payload []byte) bool {
	if secret == "" { // Doble chequeo por si acaso
		return false // No se puede validar sin secreto
	}
	if !strings.HasPrefix(ghSignature, "sha256=") {
		log.Println("ERROR: Signature format invalid (missing sha256= prefix)")
		return false
	}
	expectedSigHex := strings.TrimPrefix(ghSignature, "sha256=")
	expectedSig, err := hex.DecodeString(expectedSigHex)
	if err != nil {
		log.Printf("ERROR: Failed to decode expected signature hex: %v", err)
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload) // Usa el payload crudo
	calculatedSig := mac.Sum(nil)

	// Compara en tiempo constante para evitar ataques de temporización
	return hmac.Equal(calculatedSig, expectedSig)
}
*/