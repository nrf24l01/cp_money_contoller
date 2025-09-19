package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nrf24l01/cp_money_controller/backend/models"
	"github.com/nrf24l01/cp_money_controller/backend/schemas"
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

func (h *Handler) RefreshAccessTokenHandler(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("No refresh token found: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Refresh token not found"})
	}

	claims, err := jwtutil.ValidateToken(refreshToken.Value, []byte(h.Config.JWTRefreshSecret))
	if err != nil {
		log.Printf("Invalid refresh token: %v", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid refresh token"})
	}

	userID := claims["sub"].(string)
	username := claims["username"].(string)

	var user models.User
	if err := h.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, schemas.DefaultUnauthorizedErrorResponse)
	}

	access_token, err := jwtutil.GenerateAccessToken(userID, username, []byte(h.Config.JWTAccessSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, schemas.DefaultInternalErrorResponse)
	}

	return c.JSON(http.StatusOK, schemas.AccessTokenResponse{
		AccessToken: access_token,
	})
}