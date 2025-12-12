package telegram

// Client интерфейс для отправки сообщений в Telegram
type Client interface {
	SendMessage(token, chatID, message string) error
}
