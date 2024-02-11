package security

import "golang.org/x/crypto/bcrypt"

func HashedPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPassword(passwordHashe, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashe), []byte(password))
}
