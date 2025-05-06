package indexer

// AddToIndexRequest описывает параметры запроса для добавления к индексу
type AddToIndexRequest struct {
	Number int `query:"number" validate:"required"`
}

// AddToIndexResponse описывает успешный ответ на изменение индекса
type AddToIndexResponse struct {
	OldValue int `json:"old_value"`
	NewValue int `json:"new_value"`
	Added    int `json:"added"`
}

// ErrorResponse описывает структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
} 