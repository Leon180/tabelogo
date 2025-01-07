package utils

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

func GetAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata found in context")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return "", errors.New("no authorization token found")
	}

	token := tokens[0]
	token = strings.TrimPrefix(token, "Bearer ")

	return token, nil
}
