package model

type APIError struct {
	Error string `json:"error"`
}

func NewAPIError(msg string) *APIError {
	return &APIError{
		Error: msg,
	}
}
