package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/alesr/usrsvc/internal/users/repository"
)

// User defines domain model for a user.
type User struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// newUserDomainFromStore converts a domain model user to a storage model user.
func newUserStoreFromDomain(user *User) *repository.User {
	return &repository.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// newUserDomainFromStore converts a storage model user to a domain model user.
func newUserDomainFromStore(user *repository.User) *User {
	return &User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

const countryCodeLength = 2

type FilterParams struct {
	Country *string
}

func (f *FilterParams) normalize() {
	if f.Country != nil {
		country := *f.Country
		normalized := strings.ToUpper(strings.TrimSpace(country))
		f.Country = &normalized
	}

}

func (f *FilterParams) validate() error {
	if len(*f.Country) != countryCodeLength {
		return fmt.Errorf("could not validate country input '%s': %w", *f.Country, ErrCountryCodeInvalid)
	}
	return nil
}

type PaginationParams struct {
	Cursor string
	Limit  int
}
