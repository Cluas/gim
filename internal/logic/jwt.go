package logic

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO 从文件读取公钥
const publicKey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCnBha9sCtOxwnml3AjLdPVrdfk
7b8eFQirTZLKXW9pXv/gGhkBqIC1t4m5/jpJ7neyqCboErINMWc2tuykQ3w0TZ5Q
zKi/ehy7JcMWb6hrx8JQmDhACa0RjZ00Oo1eKObPJaavJAP1d35QL/BBl5LI9EEh
oax2pBxoSj3hwUnQYQIDAQAB
-----END PUBLIC KEY-----
`

// Member is struct of token payload
type Member struct {
	ID       int    `json:"member_id"`
	AppID    string `json:"app_id"`
	OpenID   string `json:"open_id"`
	Scene    string `json:"scene"`
	SgID     int32  `json:"sg_id"`
	Nickname string `json:"nick_name"`
	Avatar   string `json:"avatar"`
}

// MemberStdClaims is struct for more payload
type MemberStdClaims struct {
	jwt.StandardClaims
	*Member
}

// Valid is func to custom valid
func (c MemberStdClaims) Valid() (err error) {
	if c.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return errors.New("token is expired")
	}
	return
}

// JwtParseMember is func to parse and valid token
func JwtParseMember(tokenString string) (*Member, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := MemberStdClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
		if err != nil {
			return nil, fmt.Errorf("err: %v", err)
		}

		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return claims.Member, err
}
