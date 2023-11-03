package utils

import (
	"chat-be/config"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateRandomNumber(digits int) (string, error) {
	bytes := make([]byte, digits) // 3 bytes for a 6-digit number
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.New("")
	}

	// Convert bytes to an integer
	number := new(big.Int)
	number.SetBytes(bytes)

	mod, addon := 9, 1
	for i := 0; i < digits-1; i++ {
		mod *= 10
		addon *= 10
	}

	number.Mod(number, big.NewInt(int64(mod)))
	number.Add(number, big.NewInt(int64(addon)))

	// Format the integer as a 6-digit string
	code := number.String()

	fmt.Println(code, "numberGenerated")
	return code, nil
}

type Claims struct {
	Key string `json:"key"`
	jwt.RegisteredClaims
}

func GenerateToken(rawString string, expirationTime time.Time) (string, error) {
	configGet := config.GetConfig()
	jwtKey := []byte(configGet.TokenSecret)

	claims := &Claims{
		Key: rawString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)

	fmt.Println(tokenStr, "tokenStr")

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func DecodeToken(token string) (string, error) {
	configGet := config.GetConfig()
	jwtKey := []byte(configGet.TokenSecret)

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		fmt.Println(err, "errParsingToken")
		return "", err
	}

	if !tkn.Valid || claims.Key == "" {
		return "", errors.New("")
	}

	return claims.Key, nil
}

func GenerateAlphaNumeric(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	seed := time.Now().UnixNano()
	randSource := mathrand.NewSource(seed)
	rng := mathrand.New(randSource)

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}
