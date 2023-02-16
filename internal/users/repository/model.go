package repository

import "time"

// User defines storage model for a user.
type User struct {
	ID        string    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Nickname  string    `db:"nickname"`
	Password  string    `db:"password"` // This is actually a hash of the password
	Email     string    `db:"email"`
	Country   string    `db:"country"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
