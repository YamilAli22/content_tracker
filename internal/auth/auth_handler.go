package auth

import (
	"time"
	"os"
	"strconv"
	
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub uuid.UUID `json:"sub"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateJWT(id uuid.UUID, email string) (string, error) {
	var secret string = os.Getenv("JWT_SECRET_KEY")
	var expiration string = os.Getenv("JWT_EXPIRATION_HOURS")
	expiration_hours, err := strconv.Atoi(expiration) // convert string (expiration) to int (expiration_hours) (maybe a better naming ??)
	if err != nil {
 		return "", err
	}

	claims := &Claims{
			Sub: id,
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiration_hours) * time.Hour)),
					IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
 		return "", nil
	}

	return tokenString, nil
}

