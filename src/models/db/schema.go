package db

import (
	"time"
)

type User struct {
	ID        string    `db:"id primary"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *User) Fields() []string {
	fs := []string{}
	fs = append(fs,
		"id",
		"first_name",
		"last_name",
		"email",
		"password",
		"created_at",
		"updated_at",
	)
	return fs
}

func (s *User) Values() []interface{} {
	vs := []interface{}{}
	vs = append(vs,
		s.ID,
		s.FirstName,
		s.LastName,
		s.Email,
		s.Password,
		s.CreatedAt,
		s.UpdatedAt,
	)
	return vs
}
