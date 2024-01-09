package models

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
	Lg *zerolog.Logger
}

type NewUser struct {
	Email    string
	Password string
}

func (u *UserService) Create(nu *NewUser) (*User, error) {
	email := strings.ToLower(nu.Email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Lg.Error().Err(err).Msg("failed to bcrypt pw while creating user")
		return nil, fmt.Errorf("bcrypt password hash error: %w", err)
	}

	user := User{
		Email:        email,
		PasswordHash: string(hashedBytes),
	}

	row := u.DB.QueryRow(`
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2) RETURNING id`, user.Email, user.PasswordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return &user, nil
}
