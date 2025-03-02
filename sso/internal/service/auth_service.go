package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyBogomolovv/profile-manager/common/transaction"
	"github.com/SergeyBogomolovv/profile-manager/sso/internal/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, email string) (domain.User, error)
	AddAccount(ctx context.Context, userID uuid.UUID, provider domain.AccountType, password []byte) error
	AccountByID(ctx context.Context, userID uuid.UUID, provider domain.AccountType) (domain.Account, error)
}

type TokenRepo interface {
	Create(ctx context.Context, userID uuid.UUID) (string, error)
	UserID(ctx context.Context, token string) (uuid.UUID, error)
}

type authService struct {
	txManager transaction.TxManager
	users     UserRepo
	tokens    TokenRepo
	jwtSecret []byte
}

func NewAuthService(txManager transaction.TxManager, users UserRepo, tokens TokenRepo, jwtSecret []byte) *authService {
	return &authService{users: users, tokens: tokens, txManager: txManager}
}

func (s *authService) Register(ctx context.Context, email, password string) error {
	return s.txManager.Run(ctx, func(ctx context.Context) error {
		// Checks if user exists
		_, err := s.users.GetByEmail(ctx, email)
		if err == nil {
			return domain.ErrUserAlreadyExists
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			return fmt.Errorf("failed to check user exists: %w", err)
		}
		// Hash password
		hashedPassword, err := hashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		// Create user
		user, err := s.users.Create(ctx, email)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		// Create credentials account type
		if err := s.users.AddAccount(ctx, user.ID, domain.AccountTypeCredentials, hashedPassword); err != nil {
			return fmt.Errorf("failed to add account: %w", err)
		}

		// TODO: send data to rabbitmq
		return nil
	})
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

	// Create tokens
	refreshToken, err := s.tokens.Create(ctx, user.ID)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	accessToken, err := signJwt(user.ID.String(), s.jwtSecret)
	if err != nil {
		return domain.Tokens{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	userID, err := s.tokens.UserID(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			return "", domain.ErrInvalidToken
		}
		return "", fmt.Errorf("failed to get user id: %w", err)
	}
	accessToken, err := signJwt(userID.String(), s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}
	return accessToken, nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePassword(password string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
