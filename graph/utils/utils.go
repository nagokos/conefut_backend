package utils

import (
	"encoding/base64"
	"fmt"
)

func GenerateAndSetUniqueID(tableName string, id int) string {
	generateCmd := fmt.Sprintf("%s:%v", tableName, id)
	encodeString := base64.URLEncoding.EncodeToString([]byte(generateCmd))
	return encodeString
}
