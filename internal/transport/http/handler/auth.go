package handler

import (
	httpdto "moneytracker/internal/transport/http/dto"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger *zap.Logger
}

func NewAuthHandler(logger *zap.Logger) *AuthHandler {
	return &AuthHandler{logger: logger}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// get login and password from body
	var credentials httpdto.Credentials

	if err := c.BodyParser(&credentials); err != nil {
		h.logger.Error("failed to parse credential", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(httpdto.ErrorResponse{
			Error: "failed to parse credentials",
		})
	}

	// validate login and password
	expectedLogin := os.Getenv("APP_LOGIN")
	expectedPassword := os.Getenv("APP_PASSWORD")

	if credentials.Login != expectedLogin || credentials.Password != expectedPassword {
		h.logger.Error("invalid credentials",
			zap.String("received login", credentials.Login),
			zap.String("expected login", expectedLogin),
			zap.String("received password", credentials.Password),
			zap.String("expected password", expectedPassword),
		)
		return c.Status(fiber.StatusUnauthorized).JSON(httpdto.ErrorResponse{
			Error: "invalid credentials",
		})
	}

	accessToken, err := generateToken(false, time.Now().Add(httpdto.AccessTokenDuration).Unix())
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(httpdto.ErrorResponse{
			Error: "failed to generate access token",
		})
	}

	refreshToken, err := generateToken(true, time.Now().Add(httpdto.RefreshTokenDuration).Unix())
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(httpdto.ErrorResponse{
			Error: "failed to generate refresh token",
		})
	}

	// Set refresh token as an HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true, // not accessible via JS
		Secure:   false,
		SameSite: "Lax",            // or "Lax" depending on your needs
		Path:     "/",              // root path
		MaxAge:   60 * 60 * 24 * 7, // 7 days in seconds
	})

	return c.JSON(httpdto.TokenResponse{
		AccessToken: accessToken,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// verify refresh token and return new access token
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(httpdto.ErrorResponse{
			Error: "missing refresh token",
		})
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(httpdto.ErrorResponse{
			Error: "invalid refresh token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" || claims["sub"] != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(httpdto.ErrorResponse{
			Error: "invalid token",
		})
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Unix(int64(exp), 0).Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(httpdto.ErrorResponse{
			Error: "refresh token expired",
		})
	}

	newAccessToken, err := generateToken(false, time.Now().Add(httpdto.AccessTokenDuration).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(httpdto.ErrorResponse{
			Error: "failed to generate access token",
		})
	}

	return c.JSON(httpdto.TokenResponse{
		AccessToken: newAccessToken,
	})
}

func generateToken(isRefresh bool, expirationTime int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": "admin",           // subject
		"exp": expirationTime,    // expiration
		"iat": time.Now().Unix(), // issued at
	}

	if isRefresh {
		claims["type"] = "refresh" // only add this for refresh tokens
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
