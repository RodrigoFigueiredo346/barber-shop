package whatsapp

// Sender é a interface que abstrai o envio de mensagens WhatsApp.
// Troque a implementação (Evolution, Twilio) sem alterar o resto do código.
type Sender interface {
	SendMessage(phone string, message string) error
}
