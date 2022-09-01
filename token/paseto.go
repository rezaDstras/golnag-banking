package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
	"time"
)

//PasetoMaker is paseto token maker
type PasetoMaker struct {
	paseto *paseto.V2
	//encrypt payload
	symmetricKey []byte
}

//NewPasetoMaker creates a new pasetoMaker
func NewPasetoMaker(symmetrickey string) (Maker, error) {
	if len(symmetrickey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size : must be exactly %d character", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetrickey),
	}

	return maker, nil
}

func (p *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := p.paseto.Decrypt(token,p.symmetricKey,payload,nil)
	if err != nil {
		return nil , ErrInvalidToken
	}

	err = payload.Valid()
	if err !=nil {
		return nil , err
	}

	return payload , nil
}
