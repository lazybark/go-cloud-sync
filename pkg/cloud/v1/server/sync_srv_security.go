package server

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashAndSaltPassword(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("[hashAndSaltPassword]%w", err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	plainPwdBytes := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdBytes)
	if err != nil {
		return false, fmt.Errorf("[comparePasswords]%w", err)
	}

	return true, nil
}

func (s *FSWServer) createSession(log, pwd string) (user, sessionKey string) {
	return "1", "AUTH_KEY"
}

func (s *FSWServer) checkToken(hash, token string) (ok bool, err error) {
	return comparePasswords(hash, token)
}
