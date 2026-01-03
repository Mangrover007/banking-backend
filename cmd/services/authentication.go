package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/Mangrover007/banking-backend/cmd/entities"
)

type AuthService interface {
	Login(user entities.User) (string, error)
	Register(user entities.User) error
}

type auth struct {
	registeredDB map[string]entities.User // NAME -> USER
	activeUsers  map[string]string        // SID -> UID
}

var (
	ErrUserNotFound       = errors.New("you do NOT exist. register now to start existing ;)")
	ErrInvalidCredentials = errors.New("K-I-L-L yourself... imposter")
	ErrRegisterConflict   = errors.New("username taken RIP")
)

func NewAuthService(registeredDB map[string]entities.User, activeUsers map[string]string) AuthService {
	return &auth{
		registeredDB: registeredDB,
		activeUsers:  activeUsers,
	}
}

func (a *auth) Login(user entities.User) (string, error) {
	found, ok := a.registeredDB[user.Name]
	if !ok {
		return "", ErrUserNotFound
	}

	if found.Password != user.Password {
		return "", ErrInvalidCredentials
	}

	uid, _ := uuid.NewV7()
	a.activeUsers[uid.String()] = found.ID
	return uid.String(), nil
}

func (a *auth) Register(user entities.User) error {
	_, ok := a.registeredDB[user.Name]
	if ok {
		return ErrRegisterConflict
	}

	// hash the password

	// add user to DB
	// need to do this but DB would be doing this in real thing
	uid, _ := uuid.NewV7()
	user.ID = uid.String()
	a.registeredDB[user.Name] = user

	return nil
}
