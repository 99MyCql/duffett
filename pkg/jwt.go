package pkg

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyClaims JWT payload 部分
type MyClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// tokenExpireDuration 过期时间
const tokenExpireDuration = time.Hour * 2

// MySecret JWT 密钥
var MySecret []byte

// InitJwt 初始化密钥
func InitJwt() {
	MySecret = []byte(Conf.JwtSecret)
}

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	c := MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(), // 过期时间
			Issuer:    "duffett",                                  // 签发人
		},
		Username: username,
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	// .(*MyClaims) 为强制类型转换
	claims, ok := token.Claims.(*MyClaims)
	// 校验 token
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
