package auth 
import (
    "errors"
    "fastpay-backend/config"
    "fastpay-backend/pkg/hash"
    "fastpay-backend/pkg/jwt"
    "fastpay-backend/pkg/utils"
    "context"
)

type ServiceInterface interface{
	Register(req *RegisterRequest )(*AuthResponse , error)
	Login(req *LoginRequest)(*AuthResponse, error)

}

type service struct{
	repository Repository
	cfg *config.Config
}


func NewService(repo Repository , cfg *config.Config) ServiceInterface{
    return &service{
        repository: repo, 
        cfg:        cfg,
    }
}

func(s *service) Register(req *RegisterRequest )(*AuthResponse , error){
	existingUser , err := s.repository.GetUserByPhone(context.Background() , req.PhoneNumber)
	if err !=nil {return nil, err} 
	if existingUser !=nil  {return  nil, errors.New("phone number already registered") }

	hashedPassword, err := hash.HashPassword(s.cfg,req.Password)
    if err != nil {return nil, err}

	userIdentifier := utils.GenerateUserID(req.WilayaCode ,  s.cfg )

	
	user := &User{
        UserIdentifier: userIdentifier,
        PhoneNumber:    req.PhoneNumber,
        PasswordHash:   hashedPassword,
        FullName:       req.FullName,
        Role:           RoleNormal,
        IsActive:       true,
    }

	if req.Email != "" {user.Email = &req.Email}

	err = s.repository.CreateUser(context.Background() ,user) 
    if err != nil     {return nil, err}

	token, err := jwt.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
    if err != nil {return nil, err }
	return &AuthResponse{Token: token}, nil
}


func (s *service) Login(req *LoginRequest) (*AuthResponse, error) {
    user, err := s.repository.GetUserByPhone(context.Background(),req.PhoneNumber)
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

    token, err := jwt.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpiration)
    if err != nil {
        return nil, err
    }
    return &AuthResponse{Token: token}, nil
}