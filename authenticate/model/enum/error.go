package enum

type ErrorCode int

const (
	HTTPStatusOK                  ErrorCode = 200
	HTTPStatusBadRequest          ErrorCode = 400
	HTTPStatusUnauthorized        ErrorCode = 401
	HTTPStatusForbidden           ErrorCode = 403
	HTTPStatusNotFound            ErrorCode = 404
	HTTPStatusInternalServerError ErrorCode = 500
	HTTPStatusServiceUnavailable  ErrorCode = 503
)

const (
	RedisNilError ErrorCode = 1000
)

const (
	UserNotExistsError     ErrorCode = 2000
	UserExistsButNotActive ErrorCode = 2001
)

const (
	LoginUserPasswordNotMatchError ErrorCode = 3000
	SessionBlockedError            ErrorCode = 3001
	SessionExpiredError            ErrorCode = 3002
)

const (
	PlaceNotExistsError ErrorCode = 4000
)
