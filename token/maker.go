package token

import "time"

type Maker interface {

	//CreateToken creates a new token for specefic username and duration
	CreateToken(username string , duration time.Duration) (string , *Payload ,error)

	//VerfiyToken checks if token is valid or not
	VerifyToken (token string) (*Payload , error)
}
