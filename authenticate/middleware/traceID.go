package middleware

import (
	"authenticate/utility"
	"context"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type TraceIDMiddleware struct {
}

func NewTraceIDMiddleware() *TraceIDMiddleware {
	return &TraceIDMiddleware{}
}

func (a *TraceIDMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := utility.GenDefaultUUID()
		ctx.Set(string(traceIDKey), traceID)
	}
}

func GetTraceID(ctx *gin.Context) (string, bool) {
	value, exists := ctx.Get(string(traceIDKey))
	if !exists {
		return "", false
	}
	traceID, ok := value.(string)
	return traceID, ok
}

type TraceIDInterceptor struct {
}

func NewTraceIDInterceptor() *TraceIDInterceptor {
	return &TraceIDInterceptor{}
}

func (interceptor *TraceIDInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		traceID := utility.GenDefaultUUID()
		ctx = context.WithValue(ctx, traceIDKey, traceID)
		return handler(ctx, req)
	}
}

func GRPCGetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(traceIDKey).(string)
	return traceID, ok
}
