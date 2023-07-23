package apperror

import (
	"encoding/json"
)

var (
	ErrEndFound = NewErrorHandler(nil, "not found", "", "")
)

type ErrorHandler struct {
	Err              error  `json:"_"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *ErrorHandler) Error() string {
	return e.Message
}

func (e *ErrorHandler) Unwrap() error {
	return e.Err
}

func (e *ErrorHandler) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewErrorHandler(err error, message, developerMessage, code string) *ErrorHandler {
	return &ErrorHandler{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}
