package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PayLoad struct {
	UserID  string
	IsAdmin bool
	jwt.StandardClaims
}

func GenerateJWT(userID string, isAdmin bool, secret []byte) (string, error) {
	expiresAt := time.Now().Add(48 * time.Hour)
	jwtClaims := &PayLoad{
		UserID:  userID,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(cookie string, secret []byte) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(cookie, &PayLoad{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid token")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, fmt.Errorf("token is not valid or it is empty")
	}
	claims, ok := token.Claims.(*PayLoad)
	if !ok {
		return nil, fmt.Errorf("cannot parse the claims")
	}
	cred := map[string]interface{}{
		"userID": claims.UserID,
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token expired, please login again")
	}
	return cred, nil
}
