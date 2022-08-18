package helpers

import "github.com/emersion/go-bcrypt"

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return string(hash), err
	}

	return string(hash), nil
}

func CheckHashPassword(password, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
