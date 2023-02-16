/*
I could have done a more precise validation, especially
for the first, last and nickname, but I think this is enough.

Ideally, I would have used a library like https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/validator
for validating the probuf messages, but I wanted to keep this project as simple as possible.
*/
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

	// I don't think this is the best way to validate a name, but it's good enough for this project.
	// I think go languages package have a better way to do this, or some library would do it for me.
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

	// I don't think this is the best way to validate an email, but it's good enough for this project.
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrEmailFormat
	}
	return nil
}

// I hope this is not too cumbersome. I wanted to make sure that the password is somewhat secure.
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
