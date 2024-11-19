package webhook

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type HttpWebhookHandler struct {
	Client     *http.Client
	WebhookUrl string
}

type WebhookRequest struct {
	Round string `json:"round"`
	Type  string `json:"type"`
}

func (h HttpWebhookHandler) SendEvent(req WebhookRequest) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		slog.Error("error encoding req", "err", err)
		os.Exit(1)
	}
	resp, err := h.Client.Post(h.WebhookUrl, "application/json", &buf)
	if err != nil {
		slog.Error("error getting url", "err", err)
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("error status code", "code", resp.StatusCode)
	}

}
