package models

import (
	"database/sql"
	"time"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleUser  UserRole = "USER"
)

type User struct {
	Id        string       `db:"id"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Role      UserRole     `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

const (
	ConstraintsUserEmailUnique string = "email_unique"
)
