package user

import (
	"context"
	"encoding/json"
	"errors"
	"fastpay-backend/config"
	email "fastpay-backend/pkg/emails"
	"fastpay-backend/pkg/hash"
	"fastpay-backend/pkg/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
    GetProfile(ctx context.Context, userID string) (*ProfileResponse, error)
    InitiatePasswordChange(ctx context.Context, userID string, req *ChangePasswordRequest) (*ActionResponse, error)
    VerifyAndChangePassword(ctx context.Context, req *VerifyPasswordRequest) (*ActionResponse, error)
}

type service struct {
    repo   Repository
    redis  *redis.Client
    mailer *email.Mailer
}

func NewService(repo Repository, redis *redis.Client, mailer *email.Mailer) Service {
    return &service{repo: repo, redis: redis, mailer: mailer}
}

func (s *service) GetProfile(ctx context.Context, userID string) (*ProfileResponse, error) {
    user, err := s.repo.GetUserByID(ctx, userID)
    if err != nil {
        return nil, err
    }

    email := ""
    if user.Email != nil {
        email = *user.Email
    }

    return &ProfileResponse{
        ID:             user.ID,
        UserIdentifier: user.UserIdentifier,
        PhoneNumber:    user.PhoneNumber,
        Email:          email,
        FullName:       user.FullName,
        Role:           user.Role,
    }, nil
}

 func (s *service) InitiatePasswordChange(ctx context.Context, userID string, req *ChangePasswordRequest) (*ActionResponse, error) {
    // 1. Get User
    user, err := s.repo.GetUserByID(ctx, userID)
    if err != nil {
        return nil, errors.New("user not found")
    }

     if !hash.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
        return nil, errors.New("incorrect old password")
    }

    verificationToken ,err := utils.GenerateSecureToken()
    pin := utils.GeneratePIN()

    pendingData := map[string]interface{}{
        "user_id":      userID,
        "new_password": req.NewPassword, 
		
    }
    dataJSON, _ := json.Marshal(pendingData)
    
    redisKey := "pending:change_pass:" + verificationToken
    s.redis.Set(ctx, redisKey, dataJSON, 10*time.Minute)

    
    emailAddr := ""
    if user.Email != nil {
        emailAddr = *user.Email
    }
    
    subject := "FastPay Security: Password Change PIN"
    body := fmt.Sprintf("Your verification PIN is: %s. Do not share this code.", pin)
    go s.mailer.Send(emailAddr, subject, body)

    fmt.Printf("DEBUG PIN for password change: %s\n", pin) 


    return &ActionResponse{
        Status:            "verification_required",
        Message:           "A PIN has been sent to your email",
        VerificationToken: verificationToken,
    }, nil
}


func (s *service) VerifyAndChangePassword(ctx context.Context, req *VerifyPasswordRequest) (*ActionResponse, error) {
    redisKey := "pending:change_pass:" + req.VerificationToken
    val, err := s.redis.Get(ctx, redisKey).Result()
    if err != nil {
        return nil, errors.New("invalid or expired verification token")
    }

    var data map[string]interface{}
    json.Unmarshal([]byte(val), &data)

   
    pinKey := "pin:change_pass:" + req.VerificationToken
    
   
    userID := data["user_id"].(string)
    newPass := data["new_password"].(string)

   
    newHash, err := hash.HashPassword(ctx *config.Config ,newPass)
    if err != nil {
        return nil, err
    }

    
    err = s.repo.UpdatePassword(ctx, userID, newHash)
    if err != nil {
        return nil, err
    }

   
    s.redis.Del(ctx, redisKey)

    return &ActionResponse{
        Status:  "success",
        Message: "Password updated successfully",
    }, nil
}