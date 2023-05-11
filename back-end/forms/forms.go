package forms

import (
	"fmt"
	"regexp"
	"social-network/back-end/models"
	"strings"
)

func ValidateRegistrationForm(user models.User) error {
	if EmptyFields(user.Email, user.Password, user.PasswordConfirm, user.FirstName, user.LastName, user.DateOfBirth.Format("2006-01-02")) {
		return fmt.Errorf("missing fields")
	}
	if !IsDuplicate(user.Password, user.PasswordConfirm) {
		return fmt.Errorf("passwords don't match")
	}
	if !ValidLength(user.Password, 5, 15) {
		return fmt.Errorf("please enter password between 5-15 characters")
	}
	if !IsEmail(user.Email) {
		return fmt.Errorf("not valid email")
	}
	return nil
}

func EmptyFields(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func IsDuplicate(value string, value2 string) bool {
	return value == value2
}

func IsEmail(email string) bool {
	re, err := regexp.Compile(`^\S+@\S+\.\S+$`)
	if err != nil {
		fmt.Println(err)
	}
	if !re.MatchString(email) {
		return false
	}
	return true
}

func ValidLength(value string, min int, max int) bool {
	if len(value) < min || max <= len(value) {
		return false
	}
	return true
}
