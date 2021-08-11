package account_service

import (
	"crypto/md5"
	"fmt"
	"sort"
	"time"
)

func getSign(scheme *Scheme, params map[string]interface{}) string {
	m := make(map[string]interface{})
	m["appKey"] = scheme.AppKey
	m["timestamp"] = MillUnix()
	if params != nil {
		for k, v := range params {
			m[k] = v
		}
	}

	var sign = ""
	var strs []string
	for k := range m {
		strs = append(strs, k)
	}

	sort.Strings(strs)
	for _, k := range strs {
		v := fmt.Sprintf("%v", m[k])
		sign = (sign + k + v)
	}

	sign = scheme.AppSecret + sign + scheme.AppSecret
	return Xmd5(sign)
}

func Xmd5(value string) string {
	bytes := md5.Sum([]byte(value))
	format := fmt.Sprintf("%x", bytes) //将[]byte转成16进制

	return format
}

func MillUnix() int64 {
	return time.Now().UnixNano()/1e6
}
