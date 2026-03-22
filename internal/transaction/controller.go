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
func (ctrl *Controller) GetTransactionHistory(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }


    resp, err := ctrl.service.GetHistory(c.Request.Context(), userID.(string))
    if err != nil {
        c.JSON(500, gin.H{"error": "could not fetch history"})
        return
    }

    c.JSON(200, resp)
}