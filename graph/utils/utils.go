package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/nagokos/connefut_backend/logger"
)

func GenerateAndSetUniqueID(tableName string, id int) string {
	generateCmd := fmt.Sprintf("%s:%v", tableName, id)
	encodeString := base64.URLEncoding.EncodeToString([]byte(generateCmd))
	return encodeString
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
