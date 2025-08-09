package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/domain/models"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenStr, err := token.SignedString([]byte(app.Secret)) // not safety too
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
