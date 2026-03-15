package transaction

import (
	"context"
	"encoding/json"
	"errors"

	// "fastpay-backend/pkg/mail"
	email "fastpay-backend/pkg/emails"
	"fastpay-backend/pkg/jwt"
	"fastpay-backend/pkg/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
    InitiateTransfer(ctx context.Context, senderID string, req *TransferRequest, idempotencyKey string) (*TransferResponse, error)
    VerifyTransfer(ctx context.Context, req *VerifyRequest) (*TransferResponse, error)
}

type service struct {
    repo   Repository
    redis  *redis.Client
    mailer *email.Mailer
}

func NewService(repo Repository, redis *redis.Client, mailer *email.Mailer) Service {
    return &service{repo: repo, redis: redis, mailer: mailer}
}

func (s *service) InitiateTransfer(ctx context.Context, senderID string, req *TransferRequest, idempotencyKey string) (*TransferResponse, error) {
    // 1. IDEMPOTENCY CHECK
    // If key exists, it means request is processed or processing
    cacheKey := "idempotency:transfer:" + idempotencyKey
    val, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        // Return cached result if exists
        var cachedResp TransferResponse
        json.Unmarshal([]byte(val), &cachedResp)
        return &cachedResp, nil
    }

    // 2. RATE LIMIT & RULES CHECK
    today := time.Now().Format("2006-01-02")
    countKey := fmt.Sprintf("tx_count:%s:%s", senderID, today)
    
    count, _ := s.redis.Get(ctx, countKey).Int()

    // Logic: Amount > 1000 OR (Count >= 3) => Requires PIN
    requiresPIN := false
    if req.Amount > 1000 {
        requiresPIN = true
    } else if count >= 3 {
        requiresPIN = true
    }

    
    if !requiresPIN {
        // DIRECT TRANSFER
        txn, err := s.repo.ExecuteTransfer(ctx, senderID, req.ReceiverID, req.Amount)
        if err != nil {
            return nil, err
        }

        // Increment counter
        s.redis.Incr(ctx, countKey)
        s.redis.Expire(ctx, countKey, 24*time.Hour)

        // Save to Idempotency Redis
        resp := &TransferResponse{
            Status:        "completed",
            TransactionID: txn.ID,
            Message:       "Transfer successful",
        }
        data, _ := json.Marshal(resp)
        s.redis.Set(ctx, cacheKey, data, 20*time.Minute)

        return resp, nil
    }

    // 4. PIN REQUIRED FLOW
    // Generate Verification Token & PIN
    verificationToken := jwt.ValidateToken() // You need this function in utils
    pin := utils.GeneratePIN()                       // You need this function in utils (4 digits)

    // Store Pending Data in Redis (Valid for 5 minutes)
    pendingData := map[string]interface{}{
        "sender_id":       senderID,
        "receiver_id":     req.ReceiverID,
        "amount":          req.Amount,
        "idempotency_key": idempotencyKey,
        "pin":             pin,
    }
    dataJSON, _ := json.Marshal(pendingData)
    pendingKey := "pending:transfer:" + verificationToken
    s.redis.Set(ctx, pendingKey, dataJSON, 5*time.Minute)

    // Send Email (Assume we fetch email from DB, here we mock it)
    // In real app: getUserByID(senderID) -> user.Email
    s.mailer.Send("sender@example.com", "FastPay Verification PIN", fmt.Sprintf("Your PIN is: %s", pin))

    resp := &TransferResponse{
        Status:            "verification_required",
        VerificationToken: verificationToken,
        Message:           "A PIN code has been sent to your email",
    }
    
    // Cache the "verification required" state for idempotency
    data, _ := json.Marshal(resp)
    s.redis.Set(ctx, cacheKey, data, 20*time.Minute)

    return resp, nil
}

func (s *service) VerifyTransfer(ctx context.Context, req *VerifyRequest) (*TransferResponse, error) {
    // 1. Get Pending Data
    pendingKey := "pending:transfer:" + req.VerificationToken
    dataStr, err := s.redis.Get(ctx, pendingKey).Result()
    if err != nil {
        return nil, errors.New("invalid or expired transaction token")
    }

    var data map[string]interface{}
    json.Unmarshal([]byte(dataStr), &data)

    // 2. Check PIN
    if data["pin"] != req.PIN {
        return nil, errors.New("invalid PIN")
    }

    // 3. Execute Transfer
    senderID := data["sender_id"].(string)
    receiverID := data["receiver_id"].(string)
    amount := data["amount"].(float64)

    txn, err := s.repo.ExecuteTransfer(ctx, senderID, receiverID, amount)
    if err != nil {
        return nil, err
    }

    // 4. Cleanup & Update Count
    s.redis.Del(ctx, pendingKey)
    
    today := time.Now().Format("2006-01-02")
    countKey := fmt.Sprintf("tx_count:%s:%s", senderID, today)
    s.redis.Incr(ctx, countKey)
    s.redis.Expire(ctx, countKey, 24*time.Hour)

    // Update Idempotency Cache to Success
    cacheKey := "idempotency:transfer:" + data["idempotency_key"].(string)
    resp := &TransferResponse{Status: "completed", TransactionID: txn.ID}
    dataResp, _ := json.Marshal(resp)
    s.redis.Set(ctx, cacheKey, dataResp, 20*time.Minute)

    // Send Success Emails
    s.mailer.Send("sender@example.com", "Transfer Success", fmt.Sprintf("You sent %.2f", amount))
    s.mailer.Send("receiver@example.com", "Money Received", fmt.Sprintf("You received %.2f", amount))

    return resp, nil
}