package auth

import (
	"context"
	"database/sql"
	"errors"
	"log"
	appErr "mekoko/internal/errors"
	"mekoko/internal/hasher"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo           *Repository
	db             *sql.DB
	tokenGenerator TokenGenerator
}

func NewService(repo *Repository, db *sql.DB, tokenGenerator TokenGenerator) *Service {
	return &Service{repo: repo, db: db, tokenGenerator: tokenGenerator}
}

func (s *Service) Register(ctx context.Context, req RegistrationRequest) (*UserAndTokens, error) {
	email := strings.TrimSpace(req.Email)

	_, err := s.repo.FindUserByEmail(ctx, email)
	if err == nil {
		return nil, appErr.ErrUserExists
	}

	if !errors.Is(err, appErr.ErrFindingUser) {
		return nil, err
	}

	password := strings.TrimSpace(req.Password)
	confirmPassword := strings.TrimSpace(req.ConfirmPassword)

	if password != confirmPassword {
		return nil, appErr.ErrPasswordMismatch
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	publicID := uuid.NewString()
	sid := uuid.NewString()

	input := &CreateUserInput{
		PublicID:     publicID,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Email:        email,
		PasswordHash: hashedPassword,
	}

	refreshToken, jti, expiresAt, err := s.tokenGenerator.GenerateRefreshToken(publicID, sid)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(publicID, sid)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := hasher.HashToken(refreshToken)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	user, err := txRepo.CreateUser(ctx, *input)
	if err != nil {
		return nil, err
	}

	if err := txRepo.StoreRefreshToken(ctx, user.ID, sid, hashedRefreshToken, jti, expiresAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	uat := UserAndTokens{
		User:   *user,
		Tokens: tokens,
	}

	return &uat, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*Tokens, error) {
	user, err := s.repo.FindUserByEmail(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		log.Printf("Could not find user by email: %s\n", err)
		return nil, appErr.ErrInvalidCredentials
	}

	if err := ComparePassword(user.PasswordHash, strings.TrimSpace(req.Password)); err != nil {
		log.Printf("Error occured while comparing passwords: %s\n", err)
		return nil, appErr.ErrInvalidCredentials
	}

	sid := uuid.NewString()

	refreshToken, jti, expiresAt, err := s.tokenGenerator.GenerateRefreshToken(user.UUID, sid)
	hashedRefreshToken := hasher.HashToken(refreshToken)

	accessToken, err := s.tokenGenerator.GenerateAccessToken(user.UUID, sid)
	if err != nil {
		return nil, err
	}

	if err := s.repo.StoreRefreshToken(ctx, user.ID, sid, hashedRefreshToken, jti, expiresAt); err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return &tokens, nil
}

func (s *Service) Logout(ctx context.Context, sid string) error {
	return s.repo.RevokeRefreshToken(ctx, sid)
}
