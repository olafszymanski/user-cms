package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/olafszymanski/user-cms/users"
)

type JWTUserClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(user *users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTUserClaim{
		Username: *user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 120).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("could not generate jwt token, error: %w", err)
	}
	return signedToken, nil
}

func ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("could not check jwt signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
}

func ParseToken(tokenString string) (*string, error) {
	token, err := ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("could not validate jwt token, error: %w", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return &username, nil
	}
	return nil, fmt.Errorf("invalid jwt token, error: %w", err)
}
