package models

type Todo struct {
	ID    int64  `json:"id" db:"id"`
	Value string `json:"value" db:"value"`
	Date  string `json:"date" db:"date"`
}

type CreateTodoRequest struct {
	Value string `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
