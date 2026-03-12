package auth

import (
    "context"
    "errors"
    "time"

    "fastpay-backend/config"
    "fastpay-backend/pkg/hash"
    "fastpay-backend/pkg/utils"
)

type ServiceInterface interface {
    Register(req *RegisterRequest) (*AuthResponse, error)
    Login(req *LoginRequest) (*AuthResponse, error)
}

type service struct {
    repository Repository
    cfg        *config.Config
}

func NewService(repo Repository, cfg *config.Config) ServiceInterface {
    return &service{
        repository: repo,
        cfg:        cfg,
    }
}

func (s *service) createSession(ctx context.Context, userID string) (string, error) {
    rawToken, err := utils.GenerateSecureToken()
    if err != nil {
        return "", err
    }

    tokenHash := utils.HashToken(rawToken)

    session := &Session{
        UserID:       userID,
        TokenHash:    tokenHash,
        ExpiresAt:    time.Now().Add(48 * time.Hour), // 2 Days expiration
        LastActivity: time.Now(),
    }

    if err := s.repository.CreateSession(ctx, session); err != nil {
        return "", err
    }

    return rawToken, nil
}

func (s *service) Register(req *RegisterRequest) (*AuthResponse, error) {
    ctx := context.Background()

    existingUser, err := s.repository.GetUserByPhone(ctx, req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    if existingUser != nil {
        return nil, errors.New("phone number already registered")
    }

    hashedPassword, err := hash.HashPassword(&config.Config{} ,    req.Password)
    if err != nil {
        return nil, err
    }

    userIdentifier := utils.GenerateUserID(req.WilayaCode, s.cfg)

    user := &User{
        UserIdentifier: userIdentifier,
        PhoneNumber:    req.PhoneNumber,
        PasswordHash:   hashedPassword,
        FullName:       req.FullName,
        Role:           RoleNormal,
        IsActive:       true,
    }

    if req.Email != "" {
        user.Email = &req.Email
    }

    err = s.repository.CreateUser(ctx, user)
    if err != nil {
        return nil, err
    }

    token, err := s.createSession(ctx, user.ID)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{Token: token}, nil
}

func (s *service) Login(req *LoginRequest) (*AuthResponse, error) {
    ctx := context.Background()

    user, err := s.repository.GetUserByPhone(ctx, req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }

    if !user.IsActive {
        return nil, errors.New("account is disabled")
    }

    token, err := s.createSession(ctx, user.ID)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{Token: token}, nil
}