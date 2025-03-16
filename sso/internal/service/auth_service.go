package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/common/api/events"
	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Broker interface {
	PublishUserRegister(user events.UserRegister) error
}

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, email string) (domain.User, error)
	AddAccount(ctx context.Context, userID uuid.UUID, provider domain.AccountType, password []byte) (domain.Account, error)
	AccountByID(ctx context.Context, userID uuid.UUID, provider domain.AccountType) (domain.Account, error)
}

type TokenRepo interface {
	Create(ctx context.Context, userID uuid.UUID) (string, error)
	UserID(ctx context.Context, token string) (uuid.UUID, error)
	Revoke(ctx context.Context, token string) error
}

type authService struct {
	txManager transaction.TxManager
	users     UserRepo
	tokens    TokenRepo
	broker    Broker
	jwtSecret []byte
}

func NewAuthService(broker Broker, txManager transaction.TxManager, users UserRepo, tokens TokenRepo, jwtSecret []byte) *authService {
	return &authService{users: users, tokens: tokens, txManager: txManager, jwtSecret: jwtSecret, broker: broker}
}

func (s *authService) Register(ctx context.Context, email, password string) (string, error) {
	var userID uuid.UUID
	err := s.txManager.Run(ctx, func(ctx context.Context) error {
		user, added, err := s.ensureUser(ctx, email)
		if err != nil {
			return fmt.Errorf("failed to ensure user: %w", err)
		}
		userID = user.ID
		// Checks if credentials account type not exists
		_, err = s.users.AccountByID(ctx, userID, domain.AccountTypeCredentials)
		if err == nil {
			return domain.ErrUserAlreadyExists
		}
		if !errors.Is(err, domain.ErrAccountNotFound) {
			return fmt.Errorf("failed to get account: %w", err)
		}
		// Hash password
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		// Create credentials account type
		_, err = s.users.AddAccount(ctx, userID, domain.AccountTypeCredentials, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed to add account: %w", err)
		}

		if !added {
			return nil
		}
		// Publish user register
		return s.broker.PublishUserRegister(events.UserRegister{
			ID:    userID.String(),
			Email: email,
		})
	})
	if err != nil {
		return "", err
	}
	return userID.String(), nil
}

func (s *authService) Login(ctx context.Context, email, password string) (domain.Tokens, error) {
	// Get user
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.Tokens{}, domain.ErrInvalidCredentials
		}
		return domain.Tokens{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Get account
	account, err := s.users.AccountByID(ctx, user.ID, domain.AccountTypeCredentials)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return domain.Tokens{}, domain.ErrInvalidCredentials
		}
		return domain.Tokens{}, fmt.Errorf("failed to get account: %w", err)
	}

	// Compare password
	if err := comparePassword(password, account.Password); err != nil {
		return domain.Tokens{}, domain.ErrInvalidCredentials
	}

	return s.createTokens(ctx, user.ID)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.tokens.UserID(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user id: %w", err)
	}
	accessToken, err := s.signJwt(userID)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return accessToken, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.tokens.Revoke(ctx, refreshToken); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePassword(password string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
