package middleware

import (
	"bytes"
	"context"
	"google-map/model/enum"
	"google-map/utility"
	"io"
	"slices"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type TraceIDMiddleware struct {
}

func NewTraceIDMiddleware() *TraceIDMiddleware {
	return &TraceIDMiddleware{}
}

func (a *TraceIDMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventID := utility.GenDefaultUUID()
		c.Set(enum.TraceIDKey.ToString(), eventID)
		if contextType := c.Request.Header.Get(enum.RequestHeaderContentType.ToString()); slices.Contains(enum.ContextTypeGroupDefault.GetSlice().ToStringSlice(), contextType) {
			buf, err := io.ReadAll(c.Request.Body)
			if err != nil {
				utility.SugarLogger.Error(err)
				c.Next()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
			utility.SugarLogger.Infof("Http Request: %+v, EventID: %s, Body: %s", c.Request, eventID, string(buf))
		} else {
			utility.SugarLogger.Infof("Http Request: %+v, EventID: %s", c.Request, eventID)
		}
		c.Next()
	}
}

func GetTraceID(ctx *gin.Context) (string, bool) {
	value, exists := ctx.Get(enum.TraceIDKey.ToString())
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
		ctx = context.WithValue(ctx, enum.TraceIDKey, traceID)
		utility.SugarLogger.Infof("Http Request: %+v, EventID: %s", req, traceID)
		return handler(ctx, req)
	}
}

func GRPCGetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(enum.TraceIDKey).(string)
	return traceID, ok
}
