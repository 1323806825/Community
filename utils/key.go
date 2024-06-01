package utils

import (
	"fmt"
	"strings"
)

func GeneralRedisKey(originKey string, replacePart string, replaceValue string, prefix interface{}) string {
	if _, ok := prefix.(string); !ok {
		prefix = ""
	}
	return fmt.Sprintf("%s%s", prefix, strings.ReplaceAll(originKey, replacePart, replaceValue))
}
