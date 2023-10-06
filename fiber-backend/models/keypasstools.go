package models

import "golang.org/x/crypto/bcrypt"


// Generate the user's hash password 
func GenerateHasKey(key string) (string,error) {
	hash,err := bcrypt.GenerateFromPassword([]byte(key),10)
	return string(hash),err
}

// Compare plain text user's password with user's hash password
func CheckHashPassword(userhash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userhash), []byte(password))
	return err
}