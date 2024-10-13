package auth

import (
	"io"
	"os"
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
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	f, err := os.Open("C:/Go-projects/consultant-microservices/services/sso/config/private_key.pem")
	if err != nil {
		return "", errors.Wrap(err, "opening private key file")
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "reading private key file")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return "", errors.Wrap(err, "parsing private key file")
	}

	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "getting token")
	}

	return t, nil
}

func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Wrap(UnknownSignMethodError, "checking sign method")
		}

		f, err := os.Open("C:/Go-projects/consultant-microservices/services/public_key.pem")
		if err != nil {
			return "", errors.Wrap(err, "opening public key file")
		}

		buf, err := io.ReadAll(f)
		if err != nil {
			return "", errors.Wrap(err, "reading public key file")
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(buf)
		if err != nil {
			return "", errors.Wrap(err, "parsing public key file")
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "parsing token")
	}

	if !t.Valid {
		return nil, errors.Wrap(InvalidTokenError, "checking if token valid")
	}

	return claims, nil
}
