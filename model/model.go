package model

import (
	"database/sql"
	"gocommerce/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
}

// Hash Password with bcrypt
func HashPassword(password string)(string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Compare Password with bcrypt
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUser(user User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func GetUserPassword(username string) (string, error) {
	var password string
	err := database.DB.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}

	return password, nil
}
