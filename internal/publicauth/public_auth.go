package publicauth

import (
	"cm/internal/utils"
	"io"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var (
	UnknownSignMethodError = errors.New("unknown sign method")
	InvalidTokenError      = errors.New("invalid token")
)

type Claims struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	IsConsultant bool   `json:"is_consultant"`
	jwt.StandardClaims
}

func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Wrap(UnknownSignMethodError, "checking sign method")
		}

		path, err := utils.GetPath("internal/publicauth/public_key.pem")
		if err != nil {
			return "", errors.Wrap(err, "getting public_key.pem path")
		}

		f, err := os.Open(path)
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
