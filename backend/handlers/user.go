package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_contoller/models"
	"github.com/nrf24l01/cp_money_contoller/schemas"
	"github.com/nrf24l01/go-web-utils/jwtutil"
)


func (h *Handler) UserLoginHandler(c echo.Context) error {
	user_data := c.Get("validatedBody").(*schemas.UserLoginRequest)

	user := &models.User{}
	if err := h.DB.Where("username = ?", user_data.Username).First(user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, schemas.DefaultUnauthorizedErrorResponse)
	}
	
	ok, err := user.CheckPassword(user_data.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	if !ok {
		return c.JSON(http.StatusUnauthorized, schemas.DefaultUnauthorizedErrorResponse)
	}

	accessToken, refreshToken, err := jwtutil.GenerateTokenPair(user.ID.String(), user.Username, []byte(h.Config.JWTAccessSecret), []byte(h.Config.JWTRefreshSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refreshToken
	cookie.HttpOnly = true
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, schemas.AccessTokenResponse{
		AccessToken: accessToken,
	})
}