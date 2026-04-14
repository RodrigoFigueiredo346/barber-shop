package whatsapp

import "fmt"

// TwilioSender é um placeholder para implementação futura com Twilio.
type TwilioSender struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewTwilioSender(accountSID, authToken, fromNumber string) *TwilioSender {
	return &TwilioSender{AccountSID: accountSID, AuthToken: authToken, FromNumber: fromNumber}
}

func (t *TwilioSender) SendMessage(phone, message string) error {
	// TODO: implementar integração com Twilio API
	fmt.Printf("[TWILIO] Enviando para %s: %s\n", phone, message)
	return nil
}
