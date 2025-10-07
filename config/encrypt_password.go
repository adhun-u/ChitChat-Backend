package config

import (
	"golang.org/x/crypto/bcrypt"
)

/*

	THIS FILE IS MAINLY USED TO ENCRYPT AND DECRYPT THE PASSWORD

*/

// Enrypting the password
func EncryptPassword(password string) (string, error) {

	bytes, byteErr := bcrypt.GenerateFromPassword([]byte(password), 4)

	return string(bytes), byteErr
}

// Decrypting the password and checking
func CheckAndDecryptPassword(encryptedPassword string, password string) bool {

	isCorrect := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))

	return isCorrect == nil
}
