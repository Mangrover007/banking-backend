package services

import (
	"context"
	"errors"

	"github.com/Mangrover007/banking-backend/internals/repository"
	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	Register(ctx context.Context, user repository.RegisterUserParams) (repository.User, error)
	Login(ctx context.Context, phone, email, password string) (string, error)
	Logout(ctx context.Context, sid string) error
}

var (
	ErrUserNotFound       = errors.New("User is not in Database")
	ErrPhoneIsRegistered  = errors.New("Phone number is registered")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)

type service struct {
	query *repository.Queries
}

func NewAuthService(query *repository.Queries) AuthService {
	return &service{
		query: query,
	}
}

func (s *service) Register(ctx context.Context, user repository.RegisterUserParams) (repository.User, error) {
	// is phone number already registered
	_, err := s.query.FindUserByPhone(ctx, user.PhoneNumber)
	if !errors.Is(err, pgx.ErrNoRows) {
		return repository.User{}, ErrPhoneIsRegistered
	}
	// successfully register user
	newUser, err := s.query.RegisterUser(ctx, user)
	if err != nil {
		return repository.User{}, err
	}

	return newUser, nil
}

func (s *service) Login(ctx context.Context, phone, email, password string) (string, error) {
	// does user exist
	found, err := s.query.FindUserByPhoneOrEmail(ctx, repository.FindUserByPhoneOrEmailParams{
		PhoneNumber: phone,
		Email:       email,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	// is user password correct
	if password != found.Password {
		return "", ErrInvalidCredentials
	}

	// login user
	sessionID, err := s.query.CreateSession(ctx, found.ID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *service) Logout(ctx context.Context, sid string) error {
	// find and delete session
	_, err := s.query.DeleteSession(ctx, sid)
	if err != nil {
		return err
	}

	return nil
}
