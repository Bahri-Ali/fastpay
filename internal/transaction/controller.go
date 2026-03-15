package transaction

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

type Controller struct {
    service Service
}

func NewController(service Service) *Controller {
    return &Controller{service: service}
}

// InitTransfer handles POST /transfer
func (ctrl *Controller) InitTransfer(c *gin.Context) {
    var req TransferRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, _ := c.Get("user_id")
    idempotencyKey := c.GetHeader("Idempotency-Key")

    if idempotencyKey == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
        return
    }

    resp, err := ctrl.service.InitiateTransfer(c.Request.Context(), userID.(string), &req, idempotencyKey)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

// VerifyTransfer handles POST /transfer/verify
func (ctrl *Controller) VerifyTransfer(c *gin.Context) {
    var req VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := ctrl.service.VerifyTransfer(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}