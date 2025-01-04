package entitymodel

import (
	"authenticate/model/entitymodel"
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func ConvertRequest(ctx context.Context) entitymodel.Request {
	p, _ := peer.FromContext(ctx)
	return entitymodel.Request{
		UserAgent: func() string {
			md, ok := metadata.FromIncomingContext(ctx)
			if ok && len(md.Get("user-agent")) > 0 {
				return md.Get("user-agent")[0]
			}
			return ""
		}(),
		ClientIP: p.Addr.String(),
	}
}
