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

var (
	UserNotExistsError     = &APIErr{HTTPStatus: http.StatusConflict, ErrorCode: enum.UserNotExistsError, ErrorMessage: ErrorMessageMap[enum.UserNotExistsError]}
	UserExistsButNotActive = &APIErr{HTTPStatus: http.StatusConflict, ErrorCode: enum.UserExistsButNotActive, ErrorMessage: ErrorMessageMap[enum.UserExistsButNotActive]}
)

var (
	LoginUserPasswordNotMatchError = &APIErr{HTTPStatus: http.StatusUnauthorized, ErrorCode: enum.LoginUserPasswordNotMatchError, ErrorMessage: ErrorMessageMap[enum.LoginUserPasswordNotMatchError]}
	SessionBlockedError            = &APIErr{HTTPStatus: http.StatusUnauthorized, ErrorCode: enum.SessionBlockedError, ErrorMessage: ErrorMessageMap[enum.SessionBlockedError]}
)

var (
	PlaceNotExistsError = &APIErr{HTTPStatus: http.StatusNotFound, ErrorCode: enum.PlaceNotExistsError, ErrorMessage: ErrorMessageMap[enum.PlaceNotExistsError]}
)

var ErrorMessageMap = map[enum.ErrorCode]string{
	enum.HTTPStatusBadRequest: "Bad Request",
	enum.RedisNilError:        "Redis Nil Error",

	// user
	enum.UserNotExistsError:     "User already exists",
	enum.UserExistsButNotActive: "User exists but not active",

	// login
	enum.LoginUserPasswordNotMatchError: "Password not match",
	enum.SessionBlockedError:            "Session blocked",

	// place
	enum.PlaceNotExistsError: "Place not exists",
}
