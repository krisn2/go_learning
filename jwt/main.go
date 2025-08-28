package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id   uint
	Name string
}

type Claims struct {
	UserId uint
	jwt.RegisteredClaims
}

var jwtkey = []byte("krisn")

func GenerateJWT(uid uint) (string, error) {

	claims := &Claims{
		UserId: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires 24 hours from now.
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// Token was issued at the current time.
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtkey)
}

func DecodeJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtkey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !tkn.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func main() {
	user := User{
		Id:   1,
		Name: "krisn",
	}

	signedToken, err := GenerateJWT(user.Id)

	if err != nil {
		fmt.Printf("❌ Error generating token: %v\n", err)
		return
	}
	fmt.Println("✨ Generated Token:", signedToken)

	fmt.Println("---")

	decodedClaims, err := DecodeJWT(signedToken)
	if err != nil {
		fmt.Printf("❌ Error decoding token: %v\n", err)
		return
	}

	fmt.Printf("✅ Decoded Token - User ID: %d\n", decodedClaims.UserId)
	fmt.Printf("✅ Decoded Token - Expires At: %s\n", decodedClaims.ExpiresAt.Time.Format(time.RFC3339))
	fmt.Printf("✅ Decoded Token - Issued At: %s\n", decodedClaims.IssuedAt.Time.Format(time.RFC3339))
}
