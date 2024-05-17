package auth

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/fahmifan/devkit/pkg/core"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail     = errors.New("invalid email")
	ErrPasswordTooShort = errors.New("password too short")
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	Active       bool
	PasswordHash string
	Role         Role

	core.TimestampMetadata
}

type RegisterNewUserRequest struct {
	Email         string
	PlainPassword string
}

const minPassLen = 8
const bcryptCost = 12

func RegisterNewUser(newGUID uuid.UUID, req RegisterNewUserRequest) (User, error) {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return User{}, ErrInvalidEmail
	}

	if len(strings.TrimSpace(req.PlainPassword)) < minPassLen {
		return User{}, ErrPasswordTooShort
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(req.PlainPassword), bcryptCost)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		ID:           newGUID,
		Email:        req.Email,
		PasswordHash: string(passwordHashed),
		Active:       true,
	}

	return user, nil
}
