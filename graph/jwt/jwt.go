package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nagokos/connefut_backend/logger"
)

var (
	SecretKey = []byte("secretKey")
)

//* Cookieにセットする認証トークンを生成
func GenerateToken(userID int) (string, error) {
	now := time.Now().Local()
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(time.Hour * 24).Unix(),
		"iat": now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err := token.SignedString([]byte("secretKey"))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

//* CookieのtokenをパースしてユーザーIDを取得
func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, err
	}
	claims := token.Claims.(jwt.MapClaims)
	viewerID := claims["sub"].(float64)
	return int(viewerID), nil
}
