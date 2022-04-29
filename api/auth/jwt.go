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
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

func GenerateToken(user *users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTUserClaim{
		Username: *user.Username,
		Admin:    *user.Admin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 30).Unix(),
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
		return os.Getenv("SECRET_KEY"), nil
	})
}

// func ParseToken(tokenString string) (*JWTUserClaim, error) {
// 	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 	// 		return nil, fmt.Errorf("could not check jwt signing method")
// 	// 	}
// 	// 	return os.Getenv("SECRET_KEY"), nil
// 	// })
// 	token, err := ValidateToken(tokenString)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not parse jwt token, error: %w", err)
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		userClaim := &JWTUserClaim{}
// 		mapstructure.Decode(claims, userClaim)
// 		return userClaim, nil
// 	} else {
// 		return nil, fmt.Errorf("invalid authorization token, error: %w", err)
// 	}
// }
