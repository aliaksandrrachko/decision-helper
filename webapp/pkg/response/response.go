package response

import (
	"net/http"
	"time"
)

type ErrorResponse struct {
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorResponse(status int, message string, err error) ErrorResponse {
	return ErrorResponse{status, message, err.Error(), time.Now()}
}

func NewErrorResponseByCode(status int, err error) ErrorResponse {
	return ErrorResponse{
		status,
		http.StatusText(status),
		err.Error(),
		time.Now(),
	}
}
