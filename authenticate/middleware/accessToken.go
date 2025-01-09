package middleware

import (
	"authenticate/model/entitymodel"
	"authenticate/model/enum"
	"authenticate/redisDB"
	"authenticate/repository"
	"authenticate/repository/postgresqlRepository"
	"authenticate/utility"
	"context"
	"strings"

	"authenticate/errors"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"
)

const (
	authorizationPayloadKey = "authorization_payload"
)

// authMethods is a map of methods that require authentication
var authMethods = map[string]bool{
	"/grpcauth.UserService/SaveFavorite":     true,
	"/grpcauth.UserService/GetUserFavorites": true,
	"/grpcauth.UserService/LogoutUser":       true,
	"/grpcauth.PlaceService/SavePlace":       true,
	"/grpcauth.PlaceService/GetPlace":        true,
}

type AuthMiddleware struct {
	sessionRedis redisDB.SessionHandler
	sessionRepo  repository.SessionHandler
}

func NewAuthMiddleware(
	sessionRedis redisDB.SessionHandler,
	sessionRepo repository.SessionHandler,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRedis: sessionRedis,
		sessionRepo:  sessionRepo,
	}
}

func NewTransactionAuthMiddleware(
	redisClient *redis.Client,
	db *gorm.DB,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRedis: redisDB.NewSessionHandler(redisClient),
		sessionRepo:  postgresqlRepository.NewSessionHandler(db),
	}
}

func (a *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get authorization header
		authHeader := ctx.GetHeader(strings.ToLower(enum.RequestHeaderAuthorization.ToString()))
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
			return
		}

		// Check bearer token format
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
			return
		}

		authType := fields[0]
		if enum.AuthorizationTypeBearer.MatchNoCase(authType) {
			ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
			return
		}

		accessToken := fields[1]

		// Try to get session from Redis first
		session, err := a.sessionRedis.GetSession(ctx, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(500, errors.HTTPStatusInternalServerError)
			return
		}

		// If session not found in Redis, try database
		if !session.IsExist() {
			sessionWithUser, err := a.sessionRepo.GetSessionWithUserByAccessToken(ctx, accessToken)
			if err != nil {
				ctx.AbortWithStatusJSON(500, errors.HTTPStatusInternalServerError)
				return
			}
			if !sessionWithUser.IsExist() {
				ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
				return
			}
			if sessionWithUser.Session.IsAccessTokenExpired() || sessionWithUser.Session.Blocked() {
				ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
				return
			}
			// If found in DB, cache it in Redis for future requests
			if err := a.sessionRedis.SetSession(ctx, sessionWithUser.AccessToken, sessionWithUser.Session, nil); err != nil {
				utility.LogWithTraceID(ctx, "Failed to cache session in Redis:", err)
			}
			if err := a.sessionRedis.SetSession(ctx, sessionWithUser.User.Email, sessionWithUser.Session, nil); err != nil {
				utility.LogWithTraceID(ctx, "Failed to cache session in Redis:", err)
			}
		}

		// Check if session is valid
		if session.IsAccessTokenExpired() || session.Blocked() {
			ctx.AbortWithStatusJSON(401, errors.HTTPStatusUnauthorized)
			return
		}

		// Store session in context for later use
		ctx.Set(authorizationPayloadKey, session)
		ctx.Next()
	}
}

func GetSession(ctx *gin.Context) (entitymodel.Session, bool) {
	value, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return entitymodel.Session{}, false
	}
	session, ok := value.(entitymodel.Session)
	return session, ok
}

func extractAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.HTTPStatusUnauthorized
	}

	values := md.Get(strings.ToLower(enum.RequestHeaderAuthorization.ToString()))
	if len(values) == 0 {
		return "", errors.HTTPStatusUnauthorized
	}

	authStr := values[0]
	fields := strings.Fields(authStr)
	if len(fields) < 2 {
		return "", errors.HTTPStatusUnauthorized
	}
	authType := fields[0]
	if !enum.AuthorizationTypeBearer.MatchNoCase(authType) {
		return "", errors.HTTPStatusUnauthorized
	}

	accessToken := fields[1]
	if accessToken == "" {
		return "", errors.HTTPStatusUnauthorized
	}

	return accessToken, nil
}

type AuthInterceptor struct {
	sessionRedis redisDB.SessionHandler
	sessionRepo  repository.SessionHandler
}

func NewAuthInterceptor(
	sessionRedis redisDB.SessionHandler,
	sessionRepo repository.SessionHandler,
) *AuthInterceptor {
	return &AuthInterceptor{
		sessionRedis: sessionRedis,
		sessionRepo:  sessionRepo,
	}
}

func NewTransactionAuthInterceptor(
	redisClient *redis.Client,
	db *gorm.DB,
) *AuthInterceptor {
	return &AuthInterceptor{
		sessionRedis: redisDB.NewSessionHandler(redisClient),
		sessionRepo:  postgresqlRepository.NewSessionHandler(db),
	}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for methods that don't require it
		if !authMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract token
		accessToken, err := extractAccessToken(ctx)
		if err != nil {
			return nil, err
		}

		// Get and validate session
		session, err := interceptor.sessionRedis.GetSession(ctx, accessToken)
		if err != nil {
			return nil, err
		}

		if !session.IsExist() {
			sessionWithUser, err := interceptor.sessionRepo.GetSessionWithUserByAccessToken(ctx, accessToken)
			if err != nil {
				return nil, err
			}
			if !sessionWithUser.IsExist() {
				return nil, errors.HTTPStatusUnauthorized
			}

			if err := interceptor.sessionRedis.SetSession(ctx, sessionWithUser.AccessToken, sessionWithUser.Session, nil); err != nil {
				utility.LogWithTraceID(ctx, "Failed to cache session in Redis:", err)
			}
			if err := interceptor.sessionRedis.SetSession(ctx, sessionWithUser.User.Email, sessionWithUser.Session, nil); err != nil {
				utility.LogWithTraceID(ctx, "Failed to cache session in Redis:", err)
			}
		}

		if session.IsAccessTokenExpired() || session.Blocked() {
			return nil, errors.HTTPStatusUnauthorized
		}

		// Add session to context
		newCtx := context.WithValue(ctx, enum.SessionKey, session)

		// Continue with the handler
		return handler(newCtx, req)
	}
}

// Helper function to get session from context
func GRPCGetSession(ctx context.Context) (entitymodel.Session, bool) {
	session, ok := ctx.Value(enum.SessionKey).(entitymodel.Session)
	return session, ok
}
