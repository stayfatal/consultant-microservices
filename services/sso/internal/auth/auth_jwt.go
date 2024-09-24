package auth

import (
	"cm/services/sso/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	UnknownSignMethodError = errors.New("unknown sign method")
	InvalidTokenError      = errors.New("invalid token")
)

type Claims struct {
	Id int `json:"id"`
	jwt.StandardClaims
}

func CreateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	t, err := token.SignedString(config.Cfg.JWTSecret)
	if err != nil {
		return "", errors.Wrap(err, "getting token")
	}

	return t, nil
}

func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(UnknownSignMethodError, "checking sign method")
		}

		return config.Cfg.JWTSecret, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "parsing token")
	}

	if !t.Valid {
		return nil, errors.Wrap(InvalidTokenError, "checking if token valid")
	}

	return claims, nil
}
