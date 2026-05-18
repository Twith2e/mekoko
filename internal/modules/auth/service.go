package auth

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"html/template"
	"log"
	appErr "mekoko/internal/errors"
	"mekoko/internal/hasher"
	"strings"
	"time"

	"github.com/google/uuid"
)

//go:embed templates/password_reset.html
var passwordReset string

type Service struct {
	repo           *Repository
	db             *sql.DB
	tokenGenerator TokenGenerator
	emailSender    EmailSender
	clientBaseURL  string
	appName        string
}

func NewService(repo *Repository, db *sql.DB, tokenGenerator TokenGenerator, emailSender EmailSender, clientBaseURL, appName string) *Service {
	return &Service{repo: repo, db: db, tokenGenerator: tokenGenerator, emailSender: emailSender, clientBaseURL: clientBaseURL, appName: appName}
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

	userPublicID := uuid.NewString()
	sid := uuid.NewString()

	input := &CreateUserInput{
		PublicID:     userPublicID,
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Email:        email,
		PasswordHash: hashedPassword,
	}

	refreshToken, jti, expiresAt, err := s.tokenGenerator.GenerateRefreshToken(userPublicID, sid)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(userPublicID, sid)
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
		ExpiresAt:    expiresAt,
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
	if err != nil {
		return nil, err
	}
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
		ExpiresAt:    expiresAt,
	}
	return &tokens, nil
}

func (s *Service) ChangePassword(ctx context.Context, userPublicID string, payload PasswordChangeRequest) error {
	currentPw := strings.TrimSpace(payload.CurrentPassword)
	newPw := strings.TrimSpace(payload.NewPassword)
	confirmNewPw := strings.TrimSpace(payload.ConfirmNewPassword)

	user, err := s.repo.FindUserByPublicID(ctx, userPublicID)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	currentPwHash := user.PasswordHash
	if err := ComparePassword(currentPwHash, currentPw); err != nil {
		log.Printf("%s", err)
		return err
	}

	if newPw != confirmNewPw {
		return appErr.ErrPasswordMismatch
	}

	newPwHash, err := HashPassword(newPw)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	if err := txRepo.UpdateUserPasswordHash(ctx, newPwHash, user.ID); err != nil {
		return err
	}

	if err := txRepo.RevokeAllSessions(ctx, user.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) ForgotPassword(ctx context.Context, payload ForgotPasswordRequest) error {
	now := time.Now().UTC()
	email := strings.TrimSpace(payload.Email)

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, appErr.ErrFindingUser) {
			return nil
		}
		return err
	}

	count, err := s.repo.GetPasswordResetAttemptCount(ctx, user.ID)
	if err != nil {
		return err
	}

	if count >= 5 {
		return appErr.ErrTooManyRequests
	}

	token, err := GenerateToken()
	if err != nil {
		return err
	}

	hashedToken := hasher.HashToken(token)
	expiresAt := now.Add(5 * time.Minute)

	url := s.clientBaseURL + "/password/reset/" + token

	if err := s.repo.RecordPasswordResetAttempt(ctx, hashedToken, user.ID, expiresAt); err != nil {
		return err
	}

	resetEmailData := ResetEmailData{
		AppName:       s.appName,
		ResetURL:      url,
		Email:         email,
		Year:          time.Now().Year(),
		ExpiryMinutes: expiresAt.Minute(),
		Name:          user.FirstName,
	}

	tmpl, err := template.New("reset").Parse(passwordReset)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, resetEmailData); err != nil {
		log.Printf("%s", err)
		return err
	}

	htmlBody := buf.String()

	if err := s.emailSender.SendEmail(ctx, email, "Password Reset", htmlBody); err != nil {
		return err
	}

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, payload ResetPasswordRequest) error {
	token := strings.TrimSpace(payload.Token)
	if token == "" {
		return appErr.ErrInvalidRequestBody
	}

	hashedToken := hasher.HashToken(token)

	user, err := s.repo.FindUserByTokenHash(ctx, hashedToken)
	if err != nil {
		log.Printf("Error occured while finding user by token hash: %s", err)
		if errors.Is(err, sql.ErrNoRows) {
			return appErr.ErrInvalidToken
		}
		return err
	}

	newPw := strings.TrimSpace(payload.NewPassword)
	if newPw == "" {
		return appErr.ErrInvalidRequestBody
	}

	confirmNewPw := strings.TrimSpace(payload.ConfirmNewPassword)
	if confirmNewPw == "" {
		return appErr.ErrInvalidRequestBody
	}

	if confirmNewPw != newPw {
		return appErr.ErrPasswordMismatch
	}

	hashedPw, err := HashPassword(newPw)
	if err != nil {
		log.Printf("Error occured while hashing password: %s", err)
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Transaction error: %s", err)
		return err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)
	if err := txRepo.UpdateUserPasswordHash(ctx, hashedPw, user.ID); err != nil {
		log.Printf("Error occured while updating password: %s", err)
		return err
	}

	if err := txRepo.RevokeAllSessions(ctx, user.ID); err != nil {
		log.Printf("Error occured while revoking all sessions: %s", err)
		return err
	}

	if err := txRepo.MarkTokenUsed(ctx, hashedToken); err != nil {
		log.Printf("Error occured while marking token as used: %s", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error occured while commiting: %s", err)
		return err
	}

	return nil
}

func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*Tokens, error) {
	hashedRefreshToken := strings.TrimSpace(hasher.HashToken(refreshToken))
	if hashedRefreshToken == "" {
		return nil, appErr.ErrInvalidSession
	}

	row, err := s.repo.FindRefreshTokenHash(ctx, hashedRefreshToken)
	if err != nil {
		return nil, appErr.ErrInvalidSession
	}

	user, err := s.repo.FindUserByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}

	sid := uuid.NewString()

	accessToken, err := s.tokenGenerator.GenerateAccessToken(user.UUID, sid)
	if err != nil {
		return nil, err
	}

	newRefreshToken, jti, expiresAt, err := s.tokenGenerator.GenerateRefreshToken(user.UUID, sid)
	if err != nil {
		return nil, err
	}
	newHashedRefreshedToken := strings.TrimSpace(hasher.HashToken(newRefreshToken))
	if newHashedRefreshedToken == "" {
		return nil, appErr.ErrRefreshingAccessToken
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	if err := txRepo.RevokeCurrentSession(ctx, row.SID); err != nil {
		return nil, err
	}

	if err := txRepo.StoreRefreshToken(ctx, user.ID, sid, newHashedRefreshedToken, jti, expiresAt); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}

	return &tokens, nil
}

func (s *Service) Logout(ctx context.Context, sid string) error {
	log.Printf("sid: %s", sid)
	return s.repo.RevokeCurrentSession(ctx, sid)
}
