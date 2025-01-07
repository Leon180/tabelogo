package entitymodel

import (
	"authenticate/model/entitymodel"
	"authenticate/model/enum"
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func ConvertRequest(ctx context.Context) entitymodel.Request {
	p, _ := peer.FromContext(ctx)
	return entitymodel.Request{
		UserAgent: func() string {
			md, ok := metadata.FromIncomingContext(ctx)
			key := strings.ToLower(enum.RequestHeaderUserAgent.ToString())
			if ok && len(md.Get(key)) > 0 {
				return md.Get(key)[0]
			}
			return ""
		}(),
		ClientIP: p.Addr.String(),
	}
}
