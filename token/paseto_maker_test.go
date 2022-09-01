package token

import (
	"github.com/rezaDastrs/banking/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPasetoMaker(t *testing.T) {
	//create newJWT
	maker , err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)


	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	//create token
	token , err := maker.CreateToken(username,duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	//verify token
	payload , err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username , payload.Username)
	require.WithinDuration(t, issuedAt , payload.IssuedAt , time.Second)
	require.WithinDuration(t, expiredAt , payload.ExpireAt , time.Second)
}

func TestExpiredPasetoToken(t *testing.T)  {
	maker , err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token , err := maker.CreateToken(util.RandomOwner(),-time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload , err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t,err, ErrorExpiredToken.Error())
	require.Nil(t, payload)
}
