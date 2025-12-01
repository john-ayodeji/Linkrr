package utils

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		go LogError(err.Error())
		return "", err
	}

	return hashedPassword, nil
}

func ComparePasswords(password, hash string) (bool, error) {
	ok, err := argon2id.ComparePasswordAndHash(password, hash); 
	if err != nil {
		go LogError(err.Error())
		return false, err
	}

	return ok, nil
}