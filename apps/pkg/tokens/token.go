package tokens

import (
	"errors"
	"fmt"
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

func ValidateToken(tokString string)(*PayLoadToken, error){
	tok, err := jwt.Parse(tokString, func(t *jwt.Token) (interface{}, error) {
		
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method : %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)

	if !ok || !tok.Valid {
		return nil , errors.New("unauthorized")
	}

	payload := claims["payload"]
	payloadToken, ok := payload.(PayLoadToken)

	if !ok {
		return nil , errors.New("invalid payload tok")
	}

	return &payloadToken, nil

}