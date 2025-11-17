package errors

type APIError struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

func New(code, msg string) APIError {
    return APIError{Error: code, Message: msg}
}