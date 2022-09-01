package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeySize = 32

//JWTMaker is json web token maker
type JWTMaker struct {
	secretKey string
}

//NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (Maker ,error)  {
	if len(secretKey)  < minSecretKeySize{
		return nil , fmt.Errorf("invalid key size must be at least %d character",minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey} , nil
}

func (J *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload , err := NewPayload(username , duration)
	if err !=nil {
		return "", err
	}

				//first arg is jwt algorithem
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(J.secretKey))
}

func (J *JWTMaker) VerifyToken(token string) (*Payload, error) {
	//check algurithem is match or not
	keyFunc := func(token *jwt.Token) (interface{}, error){
		_ , ok :=token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil , ErrInvalidToken
		}
		return []byte(J.secretKey) , nil
	}

	jwtToken , err := jwt.ParseWithClaims(token , &Payload{},keyFunc)
	if err != nil{
		verr , ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner , ErrorExpiredToken)  {
			return nil , ErrorExpiredToken
		}
		return nil , ErrInvalidToken
	}
	playload , ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil , ErrInvalidToken
	}

	return playload , nil
}