package schemas

import "net/http"

type Message struct {
	Status string `json:"status"`
}

type Error struct {
	Error string `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var DefaultInternalErrorResponse = ErrorResponse{
	Message: "Internal Server Error",
	Code:    http.StatusInternalServerError,
}

var DefaultUnauthorizedErrorResponse = ErrorResponse{
	Message: "Unauthorized",
	Code:    http.StatusUnauthorized,
}

var DefaultNotFoundResponse = ErrorResponse{
	Message: "Resource Not Found",
	Code:    http.StatusNotFound,
}

var DefaultSuccessResponse = Message{
	Status: "success",
}

var DefaultBadRequestResponse = ErrorResponse{
	Message: "Bad Request",
	Code:    http.StatusBadRequest,
}