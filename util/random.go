package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

//int run every time when call unit package to get uniq data accoding to time now
func init()  {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generate  a random integer between min and max
func RandomInt(min , max int64) int64  {
	return min + rand.Int63n(max-min+1 )
}

// RandomString RandomSring generate a random string of lenght n
func RandomString(n int)  string{
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c:= alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string  {
	return RandomString(6)
}

func RandomMoney() int64  {
	return RandomInt(0,2000)
}

func RandomCurrency() string {
	currencies := []string{
		"EUR",
		"USD",
		"CAD",
	}

	k := len(currencies)

	return currencies[rand.Intn(k)]
}

func RandomEmail() string  {
	return fmt.Sprintf("%s@gmail.com",RandomString(6))
}
