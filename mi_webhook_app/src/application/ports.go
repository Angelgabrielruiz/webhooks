// File: src/application/ports.go
package application

// NotificationService define el puerto para enviar notificaciones.
// La capa de aplicación depende de esta interfaz, no de una implementación concreta.
type NotificationService interface {
	SendNotification(channelType string, payload DiscordPayload) error
}

// WebhookProcessor define el puerto para la lógica central de la aplicación (casos de uso).
// Adaptadores controladores (como handlers HTTP) llamarán métodos en esta interfaz.
type WebhookProcessor interface {
	ProcessPullRequestEvent(payload []byte) error
	ProcessWorkflowRunEvent(payload []byte) error
}

// --- Data Transfer Object (DTO) para Notificaciones ---
// Estructuras movidas aquí ya que son parte del contrato del puerto NotificationService.

type DiscordPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

type DiscordEmbed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Color       int            `json:"color,omitempty"` // Decimal color code
	Fields      []DiscordField `json:"fields,omitempty"`
	Footer      *DiscordFooter `json:"footer,omitempty"`
	Timestamp   string         `json:"timestamp,omitempty"` // ISO8601
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type DiscordFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}