package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string , error)  {
	HashPassword , err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("faild to hash password %w",err)
	}
	return string(HashPassword) , nil
}

func ChaeckPassword(password string , hashPassword string) error  {
	return  bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
}