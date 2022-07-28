package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nagokos/connefut_backend/logger"
)

func GenerateUniqueID(tableName string, id int) string {
	generateCmd := fmt.Sprintf("%s:%v", tableName, id)
	encodeString := base64.URLEncoding.EncodeToString([]byte(generateCmd))
	return encodeString
}

func DecodeTableName(id string) string {
	dec, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	split := strings.Split(string(dec), ":")
	tableName := split[0]
	return tableName
}

func DecodeUniqueID(id string) int {
	dec, err := base64.URLEncoding.DecodeString(id)
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	split := strings.Split(string(dec), ":")
	i, err := strconv.Atoi(split[1])
	if err != nil {
		logger.NewLogger().Error(err.Error())
	}
	return i
}

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   r.TLS != nil,
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
