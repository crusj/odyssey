package Utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	secret []byte
}
type Claim struct {
	Id        int64 `json:"id"`
	ExpiresAt int64 `json:"exp"`
	jwt.StandardClaims
}

//设置秘钥
func (it *Token) SetSecret(secret string) {
	it.secret = []byte(secret)
}

/**
 *生成TOKEN
 */
func (it *Token) CreateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 30).Unix(),
	})
	return token.SignedString(it.secret)
}

/**
 *解密TOKEN
 */
func (it *Token) DecodeToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return it.secret, nil
	})
	if err != nil {
		return 0, nil
	}
	if claims, ok := token.Claims.(*Claim); ok && token.Valid {
		//检查过期
		if time.Now().Unix() > claims.ExpiresAt {
			return 0, errors.New("令牌过期")
		} else {
			return claims.Id, nil
		}
	} else {
		return 0, nil
	}

}
