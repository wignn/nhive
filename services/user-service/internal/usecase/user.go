package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/novelhive/user-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type UserUsecase struct {
	repo      domain.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (uc *UserUsecase) Register(input domain.RegisterInput) (*domain.User, string, string, error) {
	// Validate input
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Email) == "" || len(input.Password) < 6 {
		return nil, "", "", domain.ErrInvalidInput
	}

	// Check existing
	if exists, _ := uc.repo.ExistsByEmail(input.Email); exists {
		return nil, "", "", domain.ErrEmailExists
	}
	if exists, _ := uc.repo.ExistsByUsername(input.Username); exists {
		return nil, "", "", domain.ErrUsernameExists
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user := &domain.User{
		ID:           generateID(),
		Username:     input.Username,
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hash),
		Role:         "reader",
		CreatedAt:    time.Now(),
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, "", "", err
	}

	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (uc *UserUsecase) Login(input domain.LoginInput) (*domain.User, string, string, error) {
	user, err := uc.repo.GetByEmail(strings.ToLower(input.Email))
	if err != nil {
		return nil, "", "", domain.ErrInvalidPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, "", "", domain.ErrInvalidPassword
	}

	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (uc *UserUsecase) RefreshToken(refreshTokenStr string) (string, string, error) {
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrRefreshTokenInvalid
		}
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", domain.ErrRefreshTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", domain.ErrRefreshTokenInvalid
	}

	// Ensure this is a refresh token, not an access token
	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		return "", "", domain.ErrRefreshTokenInvalid
	}

	userID, _ := claims["sub"].(string)
	if userID == "" {
		return "", "", domain.ErrRefreshTokenInvalid
	}

	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return "", "", domain.ErrUserNotFound
	}
	newAccessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := uc.generateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (uc *UserUsecase) GetProfile(userID string) (*domain.User, error) {
	return uc.repo.GetByID(userID)
}

func (uc *UserUsecase) ListUsers(page, pageSize int) ([]*domain.User, int, error) {
	return uc.repo.ListAll(page, pageSize)
}

func (uc *UserUsecase) UpdateUserRole(userID, role string) error {
	// Security: validate allowed roles
	if role != "admin" && role != "reader" {
		return domain.ErrInvalidInput
	}
	return uc.repo.UpdateRole(userID, role)
}

func (uc *UserUsecase) ValidateToken(tokenStr string) (string, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return uc.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return "", "", domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", domain.ErrInvalidToken
	}

	// Only accept access tokens (or legacy tokens without type claim)
	tokenType, _ := claims["type"].(string)
	if tokenType != "" && tokenType != "access" {
		return "", "", domain.ErrInvalidToken
	}

	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	return userID, role, nil
}

func (uc *UserUsecase) generateAccessToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"type": "access",
		"exp":  time.Now().Add(accessTokenDuration).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}

func (uc *UserUsecase) generateRefreshToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"type": "refresh",
		"exp":  time.Now().Add(refreshTokenDuration).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(uc.jwtSecret)
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
