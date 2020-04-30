package main

import (
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRET_KEY = "jwt-example"
	ACCOUNT    = "jian_liu"
)

type Claims struct {
	Account string `json:"account"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

func main() {
	now := time.Now()
	jwtId := ACCOUNT + strconv.FormatInt(now.Unix(), 10)
	role := "Member"

	claims := Claims{
		Account: ACCOUNT,
		Role:    role,
		StandardClaims: jwt.StandardClaims{
			Audience:  ACCOUNT,
			ExpiresAt: now.Add(20 * time.Second).Unix(),
			Id:        jwtId,
			IssuedAt:  now.Unix(),
			Issuer:    "generateJwt",
			NotBefore: now.Add(10 * time.Second).Unix(),
			Subject:   ACCOUNT,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Printf("JWT singing error: %v", err)
	}

	log.Print(token)
}
