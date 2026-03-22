package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fastpay-backend/pkg/emails"
	"fastpay-backend/pkg/jwt"
	"fastpay-backend/pkg/utils"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
    InitiateTransfer(ctx context.Context, senderID string, req *TransferRequest, idempotencyKey string) (*TransferResponse, error)
    VerifyTransfer(ctx context.Context, req *VerifyRequest) (*TransferResponse, error)
    GetHistory(ctx context.Context, userID string) (*TransactionListResponse, error) 
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

    cacheKey := "idempotency:transfer:" + idempotencyKey
    val, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var cachedResp TransferResponse
        json.Unmarshal([]byte(val), &cachedResp)
        return &cachedResp, nil
    }

    today := time.Now().Format("2006-01-02")
    countKey := fmt.Sprintf("tx_count:%s:%s", senderID, today)
    
    count, _ := s.redis.Get(ctx, countKey).Int()

    requiresPIN := false
    if req.Amount > 1000 {
        requiresPIN = true
    } else if count >= 3 {
        requiresPIN = true
    }

    
    if !requiresPIN {
        txn, err := s.repo.ExecuteTransfer(ctx, senderID, req.ReceiverID, req.Amount)
        if err != nil {
            return nil, err
        }

        s.redis.Incr(ctx, countKey)
        s.redis.Expire(ctx, countKey, 24*time.Hour)

        resp := &TransferResponse{
            Status:        "completed",
            TransactionID: txn.ID,
            Message:       "Transfer successful",
        }
        data, _ := json.Marshal(resp)
        s.redis.Set(ctx, cacheKey, data, 1*time.Minute)
        return resp, nil
    }

    verificationToken := jwt.ValidateToken() 
    pin := utils.GeneratePIN()                   

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

    s.mailer.Send("sender@example.com", "FastPay Verification PIN", fmt.Sprintf("Your PIN is: %s", pin))

    resp := &TransferResponse{
        Status:            "verification_required",
        VerificationToken: verificationToken,
        Message:           "A PIN code has been sent to your email",
    }
    
    data, _ := json.Marshal(resp)
    s.redis.Set(ctx, cacheKey, data, 20*time.Minute)

    return resp, nil
}

func (s *service) VerifyTransfer(ctx context.Context, req *VerifyRequest) (*TransferResponse, error) {
    pendingKey := "pending:transfer:" + req.VerificationToken
    dataStr, err := s.redis.Get(ctx, pendingKey).Result()
    if err != nil {
        return nil, errors.New("invalid or expired transaction token")
    }

    var data map[string]interface{}
    json.Unmarshal([]byte(dataStr), &data)

    if data["pin"] != req.PIN {
        return nil, errors.New("invalid PIN")
    }

    senderID := data["sender_id"].(string)
    receiverID := data["receiver_id"].(string)
    amount := data["amount"].(float64)

    txn, err := s.repo.ExecuteTransfer(ctx, senderID, receiverID, amount)
    if err != nil {
        return nil, err
    }

    s.redis.Del(ctx, pendingKey)
    
    today := time.Now().Format("2006-01-02")
    countKey := fmt.Sprintf("tx_count:%s:%s", senderID, today)
    s.redis.Incr(ctx, countKey)
    s.redis.Expire(ctx, countKey, 24*time.Hour)

    cacheKey := "idempotency:transfer:" + data["idempotency_key"].(string)
    resp := &TransferResponse{Status: "completed", TransactionID: txn.ID}
    dataResp, _ := json.Marshal(resp)
    s.redis.Set(ctx, cacheKey, dataResp, 20*time.Minute)

    s.mailer.Send("sender@example.com", "Transfer Success", fmt.Sprintf("You sent %.2f", amount))
    s.mailer.Send("receiver@example.com", "Money Received", fmt.Sprintf("You received %.2f", amount))

    return resp, nil
}

func (s *service) GetHistory(ctx context.Context, userID string) (*TransactionListResponse, error) {
    txns, err := s.repo.GetTransactionsByUserID(ctx, userID)
    if err != nil {
        return nil, err
    }

    var response TransactionListResponse

    for _, t := range txns {
        item := TransactionItemDTO{
            ID:       t.ID,
            Amount:   t.Amount,
            Currency: "DZD",
            Status:   t.Status,
            Date:     t.CreatedAt,
        }

        if t.SenderID == userID {
            item.Type = "sent"
          
            item.Counterparty = "User: " + t.ReceiverID 
            item.Amount = -t.Amount 
        } else {
            item.Type = "received"
            item.Counterparty = "From: " + t.SenderID
        }

        response.Transactions = append(response.Transactions, item)
    }

    return &response, nil
}