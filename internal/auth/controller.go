package auth

import(
	"net/http"
    "github.com/gin-gonic/gin"
)
type Controller struct{
	authService ServiceInterface
} 

func NewController(service ServiceInterface) *Controller{
	return &Controller{
		authService:service,
	}
}

func (ctrl *Controller) Register(c *gin.Context){
	var req  RegisterRequest
	
	if err:= c.ShouldBindJSON(&req) 
	err != nil {
		c.JSON(http.StatusBadRequest , gin.H{
			"err":err.Error(),
		})
		return
	}

	resp , err :=ctrl.authService.Register(&req) 
	if err!=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":err.Error(),
		})
		return
	}

 	c.JSON(http.StatusOK,resp)

}

func (ctrl *Controller) Login(c *gin.Context){
	var req  LoginRequest
	if err := c.ShouldBindJSON(&req)
	err != nil {
		c.JSON(http.StatusBadRequest , gin.H{
			"err":err.Error(),
		})
		return
	}

    token, err := ctrl.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
	c.JSON(http.StatusOK , gin.H{
		"token":token,
	})
}