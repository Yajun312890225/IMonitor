package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

// 生成32位大写MD5
func MD5Upper(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return strings.ToUpper(hex.EncodeToString(ctx.Sum(nil)))
}

func GetSign(key1, key2 string) (timestamp string, sign string) {
	timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	if key1 == "" {
		key1 = "WISESOFT"
	}
	if key2 == "" {
		key2 = "MOBILE_IM_2019"
	}
	str := "key1=" + key1 + "&timestamp=" + timestamp + "&key2=" + key2
	sign = MD5Upper(str)
	return timestamp, sign
}
