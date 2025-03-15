package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	domain "github_wb/domain/value_objects"
)

// URL de tu webhook de Discord
const discordWebhookURL = "https://discord.com/api/webhooks/1343755913252700271/qKpobvEoUuuIlqIL0X5N3EeUjzMg6OTWjahgn0UCGAno2cPr4_dcTvtpVaR1C5W0Ucc1"

// Enviar mensaje a Discord
func sendDiscordMessage(content string) {
	if content == "" {
		log.Println("Advertencia: Intento de enviar mensaje vacío a Discord.")
		return
	}

	// Crear payload
	message := map[string]string{"content": content}
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al serializar JSON para Discord: %v", err)
		return
	}

	// Enviar la solicitud HTTP
	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta para depuración
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Printf("Error: Discord devolvió el código %d. Respuesta: %s", resp.StatusCode, string(body))
	} else {
		log.Println("Mensaje enviado a Discord correctamente.")
	}
}

// Procesar eventos de Pull Request
func ProcessPullRequest(payload []byte) int {
	var eventPayload domain.PullRequestEventPayload

	// Deserializar JSON
	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		log.Printf("Error al deserializar payload: %v", err)
		return 500
	}

	// Verificar acción del Pull Request
	if eventPayload.Action == "opened" || eventPayload.Action == "closed" {
		user := eventPayload.PullRequest.User.Login
		title := eventPayload.PullRequest.Title
		url := eventPayload.PullRequest.URL

		// Mensaje para Discord
		message := fmt.Sprintf("**Pull Request %s**\n👤 **Usuario:** %s\n📌 **Título:** %s\n🔗 **URL:** %s",
			eventPayload.Action, user, title, url)

		log.Println("Enviando mensaje a Discord...")
		sendDiscordMessage(message)
	} else {
		log.Printf("Pull Request Action no es 'opened' ni 'closed': %s", eventPayload.Action)
	}

	return 200
}
