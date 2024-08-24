package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PayLoadToken struct{
	AuthId int 
	Expired time.Time
}

var secretKey = []byte("qwerty")
func GenerateToken(tok *PayLoadToken) (string, error){
	tok.Expired = time.Now().Add(10 * 60 * time.Second)
	claims := jwt.MapClaims{
		"payload" : tok.AuthId,
		"expired" : tok.Expired,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil	
}