package middleware

var authMethods = map[string]bool{
	"/grpcauth.UserService/SaveFavorite":     true,
	"/grpcauth.UserService/GetUserFavorites": true,
	"/grpcauth.UserService/LogoutUser":       true,
	"/grpcauth.PlaceService/SavePlace":       true,
	"/grpcauth.PlaceService/GetPlace":        true,
}
