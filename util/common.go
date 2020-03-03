package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 加密函数 实现php的password_hash()
func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

//验证密码 实现php的password_verify()
func PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//验证印尼手机号
func CheckIndonesiaMobile(telephone string) (string, bool) {
	telephone = strings.Replace(telephone, " ", "", -1)
	telephone = strings.Replace(telephone, "o", "0", -1)
	telephone = strings.Replace(telephone, "O", "0", -1)
	regular := `^(\+62\s?|^0)(\d{3,4}-?){2}\d{3,4}$`
	reg := regexp.MustCompile(regular)
	if strings.Index(telephone, "/") > -1 {
		arr := strings.Split(telephone, "/")
		for _, v := range arr {
			if reg.MatchString(v) {
				return "", true
			}
		}
		return "", false
	}
	return telephone, reg.MatchString(telephone)
}

//生成唯一订单号
func BuildOrderNo() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	number := r.Intn(99999)
	return "C" + fmt.Sprintf("%05d", number) + time.Now().Format("20060102150405")
}

// RemoteIp 返回远程客户端的 IP，如 192.168.1.1
func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// 获取map的key
func GetKeys(m map[string]string) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// InSlice获取一个切片并在其中查找元素。如果找到它，它将返回它的密钥，否则它将返回-1和一个错误的bool。
func InSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// 时间戳转时间字符串
func TimeFormat(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format("2006-01-02 15:04:05")
}

// 生成指定范围的随机数
func RangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max-min) + min
	return randNum
}
