package application

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	domain "github_wb/domain/value_objects"
)

const discordWebhookURL = "https://discord.com/api/webhooks/1343755913252700271/qKpobvEoUuuIlqIL0X5N3EeUjzMg6OTWjahgn0UCGAno2cPr4_dcTvtpVaR1C5W0Ucc1"

func sendToDiscord(message string) {
	payload := map[string]string{"content": message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error al serializar JSON para Discord: %v", err)
		return
	}

	_, err = http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Error al enviar mensaje a Discord: %v", err)
	}
}

func ProcessPullRequest(payload []byte) int {
	var eventPayload domain.PullRequestEventPayload

	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		log.Printf("Error al deserializar payload: %v", err)
		return 500
	}

	if eventPayload.Action == "opened" {
		user := eventPayload.PullRequest.User.Login
		title := eventPayload.PullRequest.Title
		url := eventPayload.PullRequest.URL

		message := "**Nuevo Pull Request**\n" +
			"👤 **Usuario:** " + user + "\n" +
			"📌 **Título:** " + title + "\n" +
			"🔗 **URL:** " + url

		log.Println("Enviando mensaje a Discord...")
		sendToDiscord(message)
	} else {
		log.Printf("Pull Request Action no es 'opened': %s", eventPayload.Action)
	}

	return 200
}