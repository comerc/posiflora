package telegram

// NewClient создает клиент в зависимости от режима работы
func NewClient(mockMode bool) Client {
	if mockMode {
		return NewMockClient()
	}
	return NewRealClient()
}
