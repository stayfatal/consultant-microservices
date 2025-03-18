package privateauth

import (
	"cm/libs/entities"
	"cm/libs/publicauth"
	"cm/libs/utils"

	"io"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func CreateToken(user entities.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, publicauth.Claims{
		Id:           user.Id,
		Email:        user.Email,
		IsConsultant: user.IsConsultant,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	})

	path, err := utils.GetPath("services/auth/internal/privateauth/private_key.pem")
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
