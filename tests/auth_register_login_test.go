package tests

import (
	"crypto/sha256"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	ssov1 "github.com/qliquiz/protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sso/tests/suite"
	"testing"
	"time"
)

const (
	appID          = 1
	appSecret      = "test-secret"
	passDefaultLen = 8
	deltaSeconds   = 1
)

func TestRegisterLogin_Login_HappyWay(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(
		ctx,
		&ssov1.RegisterRequest{
			Email:    email,
			Password: pass,
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUseId())

	respLogin, err := st.AuthClient.Login(
		ctx,
		&ssov1.LoginRequest{
			Email:    email,
			Password: pass,
			AppId:    appID,
		},
	)
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		signingKey := sha256.Sum256([]byte(appSecret))
		return signingKey[:], nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, respReg.GetUseId(), int64(claims["id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
