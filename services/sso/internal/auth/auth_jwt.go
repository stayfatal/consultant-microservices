package auth

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

	path, err := getFilePath("private_key.pem")
	if err != nil {
		return "", errors.Wrap(err, "getting private_key path")
	}

	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "opening private key file")
	}
	defer f.Close()

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

		path, err := getFilePath("public_key.pem")
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

func getFilePath(filename string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	dir := filepath.Dir(currentFile)

	fullPath := filepath.Join(dir, filename)

	return fullPath, nil
}
