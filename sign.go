package weixinpay

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const ChinaTimeZoneOffset = 8 * 60 * 60 //Beijing(UTC+8:00)

// NewNonceString return random string in 32 characters
func NewNonceString() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

// NewTimestampString return current timestamp
func NewTimestampString() string {
	return fmt.Sprintf("%d", time.Now().Unix()+ChinaTimeZoneOffset)
}

// Sign the parameter in form of []Param with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func Sign(params Params, key string, signType ...string) string {
	sort.Sort(params)
	preSignWithKey := params.ToQueryString() + "&key=" + key
	if len(signType) < 1 || strings.ToLower(signType[0]) == "md5" {
		return fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey)))
	} else {
		hash := hmac.New(sha256.New, []byte(key))
		hash.Write([]byte(preSignWithKey))
		return fmt.Sprintf("%X", hash.Sum(nil))
	}
}

// Verify check the sign
func Verify(in interface{}, key, correctSign string) (bool, error) {
	params, err := ToParams(in)
	if err != nil {
		return false, err
	}
	sign := Sign(params, key)
	if sign != correctSign {
		return false, fmt.Errorf("signed error: wanted %s, got %s", correctSign, sign)
	}
	return true, nil
}
