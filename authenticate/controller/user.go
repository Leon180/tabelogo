package controller

import (
	"authenticate/convert/HTTPResponse"
	"authenticate/model/entitymodel"
	"authenticate/model/requestmodel"
	"authenticate/service"
	"authenticate/utility"

	"github.com/gin-gonic/gin"
)

type RegistUserControllerHandle struct {
	registUserServiceHandler service.RegistUserServiceHandler
}

func NewRegistUserControllerHandle(registUserServiceHandler service.RegistUserServiceHandler) *RegistUserControllerHandle {
	return &RegistUserControllerHandle{registUserServiceHandler: registUserServiceHandler}
}

func (handle *RegistUserControllerHandle) RegistUser(c *gin.Context) {
	var req requestmodel.RegistUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	user, err := req.ToEntity()
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	if _, err := handle.registUserServiceHandler.RegistUser(c.Request.Context(), user); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, "success")
}

type LoginUserControllerHandle struct {
	loginUserServiceHandler service.LoginUserServiceHandler
}

func NewLoginUserControllerHandle(loginUserServiceHandler service.LoginUserServiceHandler) *LoginUserControllerHandle {
	return &LoginUserControllerHandle{loginUserServiceHandler: loginUserServiceHandler}
}

func (handle *LoginUserControllerHandle) LoginUser(c *gin.Context) {
	var req requestmodel.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	user, err := req.ToEntity()
	request := entitymodel.Request{
		UserAgent: c.Request.UserAgent(),
		ClientIP:  c.ClientIP(),
	}
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	session, err := handle.loginUserServiceHandler.LoginUser(c.Request.Context(), user, request)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, HTTPResponse.LoginResponseRef(session).ToResponse())
}

type RenewAccessTokenControllerHandle struct {
	renewAccessTokenServiceHandler service.RenewAccessTokenServiceHandler
}

func NewRenewAccessTokenControllerHandle(renewAccessTokenServiceHandler service.RenewAccessTokenServiceHandler) *RenewAccessTokenControllerHandle {
	return &RenewAccessTokenControllerHandle{renewAccessTokenServiceHandler: renewAccessTokenServiceHandler}
}

func (handle *RenewAccessTokenControllerHandle) RenewAccessToken(c *gin.Context) {
	var req requestmodel.RenewAccessTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	session, err := handle.renewAccessTokenServiceHandler.RenewAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, HTTPResponse.LoginResponseRef(session).ToResponse())
}

type SaveFavoriteControllerHandle struct {
	saveFavoriteServiceHandler service.SaveFavoriteServiceHandler
}

func NewSaveFavoriteControllerHandle(saveFavoriteServiceHandler service.SaveFavoriteServiceHandler) *SaveFavoriteControllerHandle {
	return &SaveFavoriteControllerHandle{saveFavoriteServiceHandler: saveFavoriteServiceHandler}
}

func (handle *SaveFavoriteControllerHandle) SaveFavorite(c *gin.Context) {
	var req requestmodel.SaveFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	if err := handle.saveFavoriteServiceHandler.SaveFavorite(c.Request.Context(), req.ToEntity()); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, "success")
}

type GetUserFavoritesControllerHandle struct {
	getUserFavoritesServiceHandler service.GetUserFavoritesServiceHandler
}

func NewGetUserFavoritesControllerHandle(getUserFavoritesServiceHandler service.GetUserFavoritesServiceHandler) *GetUserFavoritesControllerHandle {
	return &GetUserFavoritesControllerHandle{getUserFavoritesServiceHandler: getUserFavoritesServiceHandler}
}

func (handle *GetUserFavoritesControllerHandle) GetUserFavorites(c *gin.Context) {
	var req requestmodel.GetUserFavoritesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	favorites, err := handle.getUserFavoritesServiceHandler.GetUserFavorites(c.Request.Context(), req.ToEntity())
	if err != nil {
		utility.CommonErrorResponse(c, err, nil)
		return
	}
	utility.CommonResponse(c, HTTPResponse.UserFavoritePlacesRef(favorites).ToResponse())
}
