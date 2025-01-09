package enum

type ContextKey string

const (
	SessionKey ContextKey = "session"
	TraceIDKey ContextKey = "event_id_key"
)

func (key ContextKey) ToString() string {
	return string(key)
}
