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

func (u *UserService) Authenticate(nu *NewUser) (*User, error) {
	email := strings.ToLower(nu.Email)
	user := User{
		Email: email,
	}

	row := u.DB.QueryRow(`
		SELECT id, password_hash
		FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		u.Lg.Error().Err(err).Msg("Could not authenticate user")
		return nil, fmt.Errorf("Authenticate: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(nu.Password))
	if err != nil {
		u.Lg.Error().Err(err).Msg("password did not match")
		return nil, fmt.Errorf("password did not match: %w", err)
	}
	u.Lg.Info().Msg("User password matched")

	return &user, nil
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
