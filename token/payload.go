package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

//Diffrent types of error returned by the verfication function
var (
	ErrorExpiredToken =  errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

//Payload contains the payload data of the token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IssuedAt time.Time `json:"issued_at"`
	ExpireAt time.Time `json:"expire_at"`
}

//Valid checks if the token payload is valid or not
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpireAt) {
		return ErrorExpiredToken
	}
	return nil
}

//NewPayload creates a new token payload with specefic username and duration
func NewPayload(username string , duration time.Duration) (*Payload , error)  {
	//generate token payload
	tokenID , err := uuid.NewRandom()
	if err !=nil {
		return nil , err
	}

	payload := &Payload{
		ID:       tokenID,
		Username: username,
		IssuedAt: time.Now(),
		ExpireAt: time.Now().Add(duration),
	}

	return payload , nil
}

