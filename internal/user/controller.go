package user

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

func (ctrl *Controller) GetProfile(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    profile, err := ctrl.service.GetProfile(c.Request.Context(), userID.(string))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, profile)
}

func (ctrl *Controller) ChangePasswordInit(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    var req ChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := ctrl.service.InitiatePasswordChange(c.Request.Context(), userID.(string), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (ctrl *Controller) ChangePasswordVerify(c *gin.Context) {
    var req VerifyPasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := ctrl.service.VerifyAndChangePassword(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}