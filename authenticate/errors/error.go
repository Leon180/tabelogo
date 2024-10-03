package errors

import (
	"authenticate/model/enum"
	"net/http"
)

type APIErrorInterface interface {
	GetHttpStatus() int
	GetErrorCode() int
	Error() string
}

type APIErr struct {
	HTTPStatus   int
	ErrorCode    enum.ErrorCode
	ErrorMessage string
}

func (a APIErr) GetHttpStatus() int {
	return a.HTTPStatus
}

func (a APIErr) GetErrorCode() enum.ErrorCode {
	return a.ErrorCode
}

func (a APIErr) Error() string {
	return a.ErrorMessage
}

var (
	HTTPStatusBadRequest = &APIErr{HTTPStatus: http.StatusBadRequest, ErrorCode: enum.HTTPStatusBadRequest, ErrorMessage: ErrorMessageMap[enum.HTTPStatusBadRequest]}
	RedisNilError        = &APIErr{HTTPStatus: http.StatusBadRequest, ErrorCode: enum.RedisNilError, ErrorMessage: ErrorMessageMap[enum.RedisNilError]}
)

var ErrorMessageMap = map[enum.ErrorCode]string{
	enum.HTTPStatusBadRequest: "Bad Request",
	enum.RedisNilError:        "Redis Nil Error",
}
