package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTMiddleware(c *fiber.Ctx) error {
	fmt.Println("JWT Middleware called") // DEBUG

	tokenStr, err := extractToken(c)
	if err != nil {
		fmt.Printf("Token extraction error: %v\n", err) // DEBUG
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("Extracted token: %s\n", tokenStr) // DEBUG

	_, err = validateToken(tokenStr)
	if err != nil {
		fmt.Printf("Token validation error: %v\n", err) // DEBUG
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Println("Token validated successfully") // DEBUG
	return c.Next()
}

func extractToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization") // get header
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil // return the token
}

func validateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("token has expired")
		}
	} else {
		return nil, fmt.Errorf("exp claim missing")
	}

	if sub, ok := claims["sub"].(string); ok {
		if sub != "admin" {
			return nil, fmt.Errorf("invalid subject")
		}
	} else {
		return nil, fmt.Errorf("sub claim missing")
	}
	return token, nil
}
