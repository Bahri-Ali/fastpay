package auth 
import (
    "errors"
    "context"
    "fastpay-backend/config"
    "fastpay-backend/pkg/hash"
    "fastpay-backend/pkg/jwt"
    "fastpay-backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceInterface interface{
	Register(req *RegisterRequest )(*AuthResponse , error)
	Login(req *LoginRequest)(*AuthResponse, error)

}

type service struct{
	repository Repository
	cfg *config.Config
}


func NewService(db *pgxpool.Pool, cfg *config.Config) ServiceInterface{
	authRepo := NewRepository(db)
    return &service{
        repository: authRepo, 
        cfg:        cfg,
    }
}

func(s *service) Register(req *RegisterRequest )(*AuthResponse , error){
	existingUser , err := s.repository.GetUserByPhone(req.PhoneNumber)
	if err !=nil {return nil, err} 
	if existingUser !=nil  {return  nil, errors.New("phone number already registered") }

	hashedPassword, err := hash.HashPassword(req.Password)
    if err != nil {return nil, err}

	userIdentifier := utils.GenerateUserID(req.WilayaCode)

	
	user := &User{
        UserIdentifier: userIdentifier,
        PhoneNumber:    req.PhoneNumber,
        PasswordHash:   hashedPassword,
        FullName:       req.FullName,
        Role:           RoleNormal,
        IsActive:       true,
    }

	if req.Email != "" {user.Email = &req.Email}

	err = s.repository.CreateUser(user) 
    if err != nil     {return nil, err}

	token, err := jwt.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
    if err != nil {return nil, err }
	return &AuthResponse{Token: token}, nil
}


func (s *service) Login(req *LoginRequest) (*AuthResponse, error) {
    // 1. Get user by phone
    user, err := s.repository.GetUserByPhone(req.PhoneNumber)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, errors.New("invalid credentials")
    }

    // 2. Check password
    if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
        return nil, errors.New("invalid credentials")
    }

    // 3. Check if user is active
    if !user.IsActive {
        return nil, errors.New("account is disabled")
    }

    // 4. Generate JWT Token
    token, err := jwt.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{Token: token}, nil
}