package app

import (
	"net/mail"
	"unicode"

	apiv1 "github.com/alesr/usrsvc/proto/users/v1"
	"github.com/google/uuid"
)

const (
	minNameLength     int = 2
	maxNameLength     int = 50
	minPasswordLength int = 8
	maxPasswordLength int = 128
)

func validateCreateUserRequest(req *apiv1.CreateUserRequest) error {
	if err := validateName(req.FirstName); err != nil {
		return err
	}

	if err := validateName(req.LastName); err != nil {
		return err
	}

	if err := validateName(req.Nickname); err != nil {
		return err
	}

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if err := validateCountryCode(req.Country); err != nil {
		return err
	}
	return nil
}

func validateUpdateUserRequest(req *apiv1.UpdateUserRequest) error {
	if err := validateID(req.Id); err != nil {
		return err
	}

	if err := validateName(req.FirstName); err != nil {
		return err
	}

	if err := validateName(req.LastName); err != nil {
		return err
	}

	if err := validateName(req.Nickname); err != nil {
		return err
	}

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if err := validateCountryCode(req.Country); err != nil {
		return err
	}
	return nil
}

func validateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return ErrNameFormat
		}
	}

	if len(name) < minNameLength || len(name) > maxNameLength {
		return ErrNameLength
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrEmailFormat
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < minPasswordLength || len(password) > maxPasswordLength {
		return ErrPasswordLength
	}

	var hasNumber, hasLetter, hasSpecial bool
	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNumber = true
		}
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			hasSpecial = true
		}
	}

	if !hasNumber || !hasLetter || !hasSpecial {
		return ErrPasswordFormat
	}
	return nil
}

func validateID(id string) error {
	if id == "" {
		return ErrIDRequired
	}

	if _, err := uuid.Parse(id); err != nil {
		return ErrIDFormat
	}
	return nil
}

func validateCountryCode(country string) error {
	if country == "" {
		return ErrCountryCodeRequired
	}

	if len(country) != 2 {
		return ErrCountryCodeInvalid
	}
	return nil
}
