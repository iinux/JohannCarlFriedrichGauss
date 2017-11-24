package main

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"time"
)

type user struct {
	ID          uint
	PhoneNumber string
	Name        string
}

type userClaims struct {
	User user
	jwt.StandardClaims
}

func issueJWT() string {
	issuedAt := time.Now().Unix()
	expireAt := time.Now().Add(time.Hour * 24 * 365 * 10).Unix()

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := userClaims{
		user{},
		jwt.StandardClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expireAt,
			Issuer:    "hello.com",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte("conf.JWT_SECRET"))

	return signedToken
}

func Auth() bool {
	// Parse, validate and return a jwtToken.
	token := ""

	jwtToken, err := jwt.ParseWithClaims(token, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Prevents a known exploit
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte("conf.JWT_SECRET"), nil
	})

	if err != nil {
		fmt.Println(err)
		return false
	}

	if claims, ok := jwtToken.Claims.(*userClaims); ok && jwtToken.Valid {
		fmt.Println(claims)
		return true
	} else {
		return false
	}
}

func main()  {
	fmt.Println(issueJWT())
}
