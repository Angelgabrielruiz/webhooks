// File: src/application/webhook_service.go
package application

import (
	"encoding/json"
	"fmt"
	"log"
	"time" // Importa time

	// --- IMPORTACIÓN ACTUALIZADA (usa tu nombre de módulo) ---
	domain "mi_webhook_app/src/domain/value_objects"
	// --- FIN IMPORTACIÓN ACTUALIZADA ---
)

// webhookService implementa la interfaz WebhookProcessor.
type webhookService struct {
	// Depende del puerto NotificationService (interfaz), no de una implementación concreta.
	notifier NotificationService
}

// NewWebhookService es el constructor para webhookService.
// Recibe la implementación concreta del notificador a través de la interfaz.
func NewWebhookService(notifier NotificationService) WebhookProcessor {
	return &webhookService{
		notifier: notifier,
	}
}

// ProcessPullRequestEvent maneja eventos pull_request. Implementa WebhookProcessor.
func (s *webhookService) ProcessPullRequestEvent(payload []byte) error {
	var event domain.PullRequestEventPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("ERROR: Unmarshalling PullRequestEventPayload: %v", err)
		return fmt.Errorf("failed to unmarshal pull request payload: %w", err)
	}

	var messageEmbed DiscordEmbed
	var sendMessage bool = true // Flag para controlar el envío
	var sendErr error = nil    // Para capturar errores potenciales de notificación

	pr := event.PullRequest
	repo := event.Repository
	sender := event.Sender

	// Footer y Timestamp por defecto (pueden ser sobreescritos)
	footer := &DiscordFooter{Text: fmt.Sprintf("Triggered by %s", sender.Login)}
	timestamp := time.Now().Format(time.RFC3339) // Usa tiempo actual por defecto

	switch event.Action {
	case "opened":
		if pr.CreatedAt != nil { // Verifica si CreatedAt está disponible
			timestamp = pr.CreatedAt.Format(time.RFC3339)
		}
		messageEmbed = DiscordEmbed{
			Title:       fmt.Sprintf("🚀 New Pull Request #%d: %s", event.Number, pr.Title),
			Description: fmt.Sprintf("A new pull request was opened in [%s](%s).", repo.FullName, repo.HTMLURL),
			URL:         pr.HTMLURL,
			Color:       3447003, // Azul
			Fields: []DiscordField{
				{Name: "Author", Value: fmt.Sprintf("[%s](%s)", pr.User.Login, pr.User.HTMLURL), Inline: true},
				{Name: "Branch", Value: fmt.Sprintf("`%s` → `%s`", pr.Head.Ref, pr.Base.Ref), Inline: true},
			},
			Footer:    footer,
			Timestamp: timestamp,
		}
	case "reopened":
		if pr.UpdatedAt != nil { // Verifica si UpdatedAt está disponible
			timestamp = pr.UpdatedAt.Format(time.RFC3339)
		}
		footer.Text = fmt.Sprintf("Reopened by %s", sender.Login)
		messageEmbed = DiscordEmbed{
			Title:       fmt.Sprintf("🔄 Pull Request Reopened #%d: %s", event.Number, pr.Title),
			Description: fmt.Sprintf("Pull request reopened in [%s](%s).", repo.FullName, repo.HTMLURL),
			URL:         pr.HTMLURL,
			Color:       16776960, // Amarillo
			Fields: []DiscordField{
				{Name: "Author", Value: fmt.Sprintf("[%s](%s)", pr.User.Login, pr.User.HTMLURL), Inline: true},
				{Name: "Branch", Value: fmt.Sprintf("`%s` → `%s`", pr.Head.Ref, pr.Base.Ref), Inline: true},
			},
			Footer:    footer,
			Timestamp: timestamp,
		}
	case "ready_for_review":
		if pr.UpdatedAt != nil {
			timestamp = pr.UpdatedAt.Format(time.RFC3339)
		}
		footer.Text = fmt.Sprintf("Marked ready by %s", sender.Login)
		messageEmbed = DiscordEmbed{
			Title:       fmt.Sprintf("👀 PR Ready for Review #%d: %s", event.Number, pr.Title),
			Description: fmt.Sprintf("Pull request marked as ready for review in [%s](%s).", repo.FullName, repo.HTMLURL),
			URL:         pr.HTMLURL,
			Color:       3066993, // Verde
			Fields: []DiscordField{
				{Name: "Author", Value: fmt.Sprintf("[%s](%s)", pr.User.Login, pr.User.HTMLURL), Inline: true},
				{Name: "Branch", Value: fmt.Sprintf("`%s` → `%s`", pr.Head.Ref, pr.Base.Ref), Inline: true},
			},
			Footer:    footer,
			Timestamp: timestamp,
		}
	case "closed":
		if pr.Merged {
			if pr.MergedAt != nil { // Verifica si MergedAt está disponible
				timestamp = pr.MergedAt.Format(time.RFC3339)
			}
			footer.Text = "Merged"
			messageEmbed = DiscordEmbed{
				Title:       fmt.Sprintf("✅ Pull Request Merged #%d: %s", event.Number, pr.Title),
				Description: fmt.Sprintf("Pull request successfully merged into `%s` in [%s](%s).", pr.Base.Ref, repo.FullName, repo.HTMLURL),
				URL:         pr.HTMLURL,
				Color:       8359053, // Púrpura
				Fields: []DiscordField{
					{Name: "Author", Value: fmt.Sprintf("[%s](%s)", pr.User.Login, pr.User.HTMLURL), Inline: true},
					{Name: "Merged By", Value: fmt.Sprintf("[%s](%s)", sender.Login, sender.HTMLURL), Inline: true}, // Asume que sender es quien hizo merge
				},
				Footer:    footer,
				Timestamp: timestamp,
			}
		} else {
			log.Printf("INFO: Pull Request #%d closed without merging (Action: %s). No notification sent.", event.Number, event.Action)
			sendMessage = false
		}
	default:
		log.Printf("INFO: Unhandled Pull Request Action: %s for PR #%d. No notification sent.", event.Action, event.Number)
		sendMessage = false
	}

	if sendMessage {
		log.Printf("INFO: Sending Pull Request notification to Development channel for action: %s", event.Action)
		discordPayload := DiscordPayload{Embeds: []DiscordEmbed{messageEmbed}}
		// Usa el notificador inyectado a través del puerto de interfaz
		err := s.notifier.SendNotification("development", discordPayload)
		if err != nil {
			log.Printf("ERROR: Sending development notification: %v", err)
			// Decide cómo manejar errores de notificación. Lo retornamos aquí.
			sendErr = fmt.Errorf("failed to send development notification: %w", err)
		}
	}

	return sendErr // Retorna nil si fue exitoso o si no se envió mensaje intencionalmente
}

// ProcessWorkflowRunEvent maneja eventos workflow_run. Implementa WebhookProcessor.
func (s *webhookService) ProcessWorkflowRunEvent(payload []byte) error {
	var event domain.WorkflowRunEventPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("ERROR: Unmarshalling WorkflowRunEventPayload: %v", err)
		return fmt.Errorf("failed to unmarshal workflow run payload: %w", err)
	}

	// Solo nos interesan workflows completados
	if event.Action != "completed" {
		log.Printf("INFO: Ignoring workflow_run event with action '%s'", event.Action)
		return nil // No es un error, solo ignorado
	}

	run := event.WorkflowRun
	repo := event.Repository
	workflow := event.Workflow

	var messageEmbed DiscordEmbed
	var color int
	var statusEmoji string
	var sendMessage bool = true
	var sendErr error = nil

	switch run.Conclusion {
	case "success":
		color = 3066993
		statusEmoji = "✅"
	case "failure":
		color = 15158332
		statusEmoji = "❌"
	case "cancelled":
		color = 9807270
		statusEmoji = "⏹️"
	case "skipped":
		color = 16776960
		statusEmoji = "⏭️"
	default:
		log.Printf("INFO: Unhandled Workflow Conclusion: %s for Run ID %d. No notification sent.", run.Conclusion, run.ID)
		sendMessage = false
	}

	if sendMessage {
		// Construye descripción, menciona PR si está disponible
		description := fmt.Sprintf("Workflow **%s** completed with status: **%s**", workflow.Name, run.Conclusion)
		if len(run.PullRequests) > 0 {
			prNumber := run.PullRequests[0].Number
			prURL := fmt.Sprintf("%s/pull/%d", repo.HTMLURL, prNumber)
			description += fmt.Sprintf("\nAssociated Pull Request: [#%d](%s)", prNumber, prURL)
		}

		timestamp := run.UpdatedAt.Format(time.RFC3339) // Usa tiempo de completado

		messageEmbed = DiscordEmbed{
			Title:       fmt.Sprintf("%s Workflow Run %s: %s", statusEmoji, run.Conclusion, workflow.Name),
			Description: description,
			URL:         run.HTMLURL, // Enlace a la ejecución específica
			Color:       color,
			Fields: []DiscordField{
				{Name: "Repository", Value: fmt.Sprintf("[%s](%s)", repo.FullName, repo.HTMLURL), Inline: true},
				{Name: "Branch", Value: fmt.Sprintf("`%s`", run.HeadBranch), Inline: true},
				{Name: "Triggered By", Value: fmt.Sprintf("[%s](%s)", event.Sender.Login, event.Sender.HTMLURL), Inline: true},
				{Name: "Event", Value: run.Event, Inline: true},
				{Name: "Run ID", Value: fmt.Sprintf("[%d](%s)", run.ID, run.HTMLURL), Inline: true},
			},
			Footer:    &DiscordFooter{Text: fmt.Sprintf("Workflow: %s", workflow.Path)},
			Timestamp: timestamp,
		}

		log.Printf("INFO: Sending Workflow Run notification to Testing channel (Conclusion: %s)", run.Conclusion)
		discordPayload := DiscordPayload{Embeds: []DiscordEmbed{messageEmbed}}
		// Usa el notificador inyectado a través del puerto de interfaz
		err := s.notifier.SendNotification("testing", discordPayload)
		if err != nil {
			log.Printf("ERROR: Sending testing notification: %v", err)
			sendErr = fmt.Errorf("failed to send testing notification: %w", err)
		}
	}
	return sendErr
}