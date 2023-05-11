package main

import (
	"database/sql"
	"log"
	database "social-network/back-end/database"
	"social-network/back-end/models"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func login(email, password string) bool {
	var dbPassword string
	err := database.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Println("login() error: ", err)
	}
	if CheckPasswordHash(password, dbPassword) {
		return true
	}
	return false
}

func registerUser(user models.User) (userID int, err error) {
	hashedPw, err := HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	result, err := database.Exec("INSERT INTO users (`email`,`password`,`firstName`,`lastName`,`dateOfBirth`,`img`,`nickname`,`about`,`profilePublic`) VALUES (?,?,?,?,?,?,?,?,?)", user.Email, hashedPw, user.FirstName, user.LastName, user.DateOfBirth.Format("2006-01-02"), user.Image, user.Nickname, user.About, user.ProfileP)
	if err != nil {
		return 0, err
	}
	userIDint64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userIDint64), nil
}
