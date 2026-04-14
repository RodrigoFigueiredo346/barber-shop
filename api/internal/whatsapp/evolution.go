package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EvolutionSender struct {
	BaseURL  string
	APIKey   string
	Instance string
}

func NewEvolutionSender(baseURL, apiKey, instance string) *EvolutionSender {
	return &EvolutionSender{BaseURL: baseURL, APIKey: apiKey, Instance: instance}
}

func (e *EvolutionSender) SendMessage(phone, message string) error {
	payload := map[string]interface{}{
		"number": phone,
		"text":   message,
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/message/sendText/%s", e.BaseURL, e.Instance)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", e.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("evolution api error: status %d", resp.StatusCode)
	}
	return nil
}
