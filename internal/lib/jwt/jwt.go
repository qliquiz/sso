package jwt

import (
	"crypto/sha256"
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/domain/models"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	signingKey := sha256.Sum256([]byte(app.Secret))

	tokenStr, err := token.SignedString(signingKey[:])
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
