package middleware

type contextKey string

const (
	sessionKey contextKey = "session"
	traceIDKey contextKey = "trace_id"
)
