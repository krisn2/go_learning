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

var jwtkey = []byte("Secret")

func GenarateJWT(uid uint) (string, error) {
	claims := &Claims{
		UserId: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtkey)
}

func main() {

	user := User{
		Id:   1,
		Name: "krisn",
	}

	t, err := GenarateJWT(user.Id)

	if err != nil {
		fmt.Println("Error to generate a token")
		return
	} else {
		println("Token :", t)
	}

}
