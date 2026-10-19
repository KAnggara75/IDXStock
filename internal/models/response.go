package models

// ErrorResponse represents standard error payload
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents standard success message payload
type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
