package grpcservice

import (
	"authenticate/errors"
	convertentity "authenticate/grpc/convert/entitymodel"
	convertrequest "authenticate/grpc/convert/requestmodel"
	convertresponse "authenticate/grpc/convert/responsemodel"
	"authenticate/grpc/proto"
	"authenticate/middleware"
	"authenticate/service"
	"context"
)

func NewUserServiceServer(
	registUserServiceHandler service.RegistUserServiceHandler,
	loginUserServiceHandler service.LoginUserServiceHandler,
	renewAccessTokenServiceHandler service.RenewAccessTokenServiceHandler,
	logoutUserServiceHandler service.LogoutUserServiceHandler,
	saveFavoriteServiceHandler service.SaveFavoriteServiceHandler,
	getUserFavoritesServiceHandler service.GetUserFavoritesServiceHandler,
) *UserServiceServer {
	return &UserServiceServer{
		registUserServiceHandler:       registUserServiceHandler,
		loginUserServiceHandler:        loginUserServiceHandler,
		renewAccessTokenServiceHandler: renewAccessTokenServiceHandler,
		logoutUserServiceHandler:       logoutUserServiceHandler,
		saveFavoriteServiceHandler:     saveFavoriteServiceHandler,
		getUserFavoritesServiceHandler: getUserFavoritesServiceHandler,
	}
}

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	registUserServiceHandler       service.RegistUserServiceHandler
	loginUserServiceHandler        service.LoginUserServiceHandler
	renewAccessTokenServiceHandler service.RenewAccessTokenServiceHandler
	logoutUserServiceHandler       service.LogoutUserServiceHandler
	saveFavoriteServiceHandler     service.SaveFavoriteServiceHandler
	getUserFavoritesServiceHandler service.GetUserFavoritesServiceHandler
}

func (server *UserServiceServer) RegistUser(ctx context.Context, req *proto.RegistUserRequest) (*proto.User, error) {
	request := convertrequest.ConvertRegistUserRequest(req)
	user, err := request.ToEntity()
	if err != nil {
		return nil, err
	}
	user, err = server.registUserServiceHandler.RegistUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return convertresponse.UserEntityModel(user).ToUserProto(), nil
}

func (server *UserServiceServer) LoginUser(ctx context.Context, req *proto.LoginUserRequest) (*proto.LoginUserResponse, error) {
	request := convertrequest.ConvertLoginUserRequest(req)
	user, err := request.ToEntity()
	if err != nil {
		return nil, err
	}
	session, err := server.loginUserServiceHandler.LoginUser(ctx, user, convertentity.ConvertRequest(ctx))
	if err != nil {
		return nil, err
	}
	return convertresponse.SessionEntityModel(session).ToLoginUserResponseProto(), nil
}

func (server *UserServiceServer) RenewAccessToken(ctx context.Context, req *proto.RenewAccessTokenRequest) (*proto.LoginUserResponse, error) {
	request := convertrequest.ConvertRenewAccessTokenRequest(req)
	session, err := server.renewAccessTokenServiceHandler.RenewAccessToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}
	return convertresponse.SessionEntityModel(session).ToLoginUserResponseProto(), nil
}

func (server *UserServiceServer) LogoutUser(ctx context.Context, req *proto.CommonRequest) (*proto.CommonResponse, error) {
	session, ok := middleware.GRPCGetSession(ctx)
	if !ok {
		return nil, errors.HTTPStatusUnauthorized
	}
	if err := server.logoutUserServiceHandler.LogoutUser(ctx, session.UserID, session.AccessToken); err != nil {
		return nil, err
	}
	return &proto.CommonResponse{
		Message: "success",
	}, nil
}

func (server *UserServiceServer) SaveFavorite(ctx context.Context, req *proto.SaveFavoriteRequest) (*proto.CommonResponse, error) {
	request := convertrequest.ConvertSaveFavoriteRequest(req)
	session, ok := middleware.GRPCGetSession(ctx)
	if !ok {
		return nil, errors.HTTPStatusUnauthorized
	}
	if err := server.saveFavoriteServiceHandler.SaveFavorite(ctx, request.ToEntity(session.UserID)); err != nil {
		return nil, err
	}
	return &proto.CommonResponse{
		Message: "success",
	}, nil
}

func (server *UserServiceServer) GetUserFavorites(ctx context.Context, req *proto.GetUserFavoritesRequest) (*proto.GetUserFavoritesResponse, error) {
	request := convertrequest.ConvertGetUserFavoritesRequest(req)
	session, ok := middleware.GRPCGetSession(ctx)
	if !ok {
		return nil, errors.HTTPStatusUnauthorized
	}
	favorites, err := server.getUserFavoritesServiceHandler.GetUserFavorites(ctx, request.ToEntity(session.UserID))
	if err != nil {
		return nil, err
	}
	return convertresponse.UserFavoritePlacesEntityModel(favorites).ToGetUserFavoritesResponseProto(), nil
}
