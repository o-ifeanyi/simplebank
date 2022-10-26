package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(8)

	hashed1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed1)

	err = CompareHashAndPassword(hashed1, password)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CompareHashAndPassword(hashed1, wrongPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashed2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hashed1, hashed2)
}
