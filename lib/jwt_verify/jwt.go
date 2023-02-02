package jwt_verify

/**
 * @Author: chengming1
 * @Date: 2023/2/2 15:21
 * @Desc: json web token
 * 参考：https://blog.csdn.net/neweastsun/article/details/105919915
 */

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("cmgo")

type Claims struct {
	jwt.StandardClaims

	Username string `json:"username"`
}

func GenerateToken(username string, expireTimeConf int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireTimeConf) * time.Hour)

	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
		username,
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// 解析Tokne
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("token已过期")
			} else {
				return nil, errors.New("token无效")
			}
		}
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无效")
}
