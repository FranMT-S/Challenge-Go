package model

import "fmt"

type ResponseError struct {
	Status int
	Msg    string
	Err    string
}

func NewResponseError(status int, msg string, err string) *ResponseError {
	return &ResponseError{status, msg, err}
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf(`{
			"status":%v,
			"msg":"%v",
			"error":"%v"
		}`, err.Status, err.Msg, err.Err)
}
