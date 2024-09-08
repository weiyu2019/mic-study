package jwt_op

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mic-study/conf"
	"mic-study/log"
	"time"
)

const (
	TokenExpired     = "Token is expired"
	TokenNotValidYet = "Token not active yet"
	TokenMalformed   = "That's not even a token"
	TokenInvalid     = "That's not even a token"
)

type CustomClaims struct {
	jwt.StandardClaims
	ID          int32
	NickName    string
	AuthorityId int32
}

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{SigningKey: []byte(conf.AppConf.JWTConfig.SigningKey)}
}

func (j *JWT) GenerateJWT(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.SigningKey)
	if err != nil {
		log.Logger.Error("生成JTW错误:" + err.Error())
	}
	return tokenStr, nil
}

func (j *JWT) ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if result, ok := err.(jwt.ValidationError); ok {
			if result.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(TokenMalformed)
			} else if result.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New(TokenExpired)
			} else if result.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(TokenNotValidYet)
			} else {
				return nil, errors.New(TokenInvalid)
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New(TokenInvalid)
	} else {
		return nil, errors.New(TokenInvalid)
	}
}

func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.GenerateJWT(*claims)
	}
	return "", errors.New(TokenInvalid)
}
