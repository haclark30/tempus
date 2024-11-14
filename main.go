package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"net/http"
	"os"
)

type WebhookRequest struct {
	Round string `json:"round"`
	Type  string `json:"type"`
}

func main() {

	webhookUrl := os.Getenv("WEBHOOK_URL")
	args := os.Args[1:]
	if len(args) < 1 {
		slog.Error("not enough args")
		os.Exit(1)
	}
	webhookReq := WebhookRequest{"Work Session", "Complete"}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(webhookReq); err != nil {
		slog.Error("error encoding req", "err", err)
		os.Exit(1)
	}
	resp, err := http.Post(webhookUrl, "application/json", &buf)
	if err != nil {
		slog.Error("error getting url", "err", err)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("error status code", "code", resp.StatusCode)
	}
}
