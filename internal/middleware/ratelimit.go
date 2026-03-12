package middleware

import (
    "net/http"
    "sync"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)



type IPRateLimite struct{
	ips map[string]*rate.Limiter
	mu sync.Mutex
	r rate.Limit
	b int
}

func NewIPRateLimiter(r rate.Limit , b int) *IPRateLimite {
	return &IPRateLimite{
		ips: make(map[string]*rate.Limiter),
		r:r,
		b:b,
	}
}

func (i *IPRateLimite) GetLimiter(ip string) *rate.Limiter{
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter , exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r , i.b)
		i.ips[ip] = limiter
	} 
	return  limiter

}

func RateLimit() gin.HandlerFunc {
	limiter := NewIPRateLimiter(1 , 2)
	return func (c *gin.Context){
		ip:=c.ClientIP()

		if !limiter.GetLimiter(ip).Allow(){
			c.JSON(http.StatusTooManyRequests , gin.H{
				"err":"to many request !",
			})
			c.Abort()
			return
		}
		c.Next()

		if len(limiter.ips)>1000{
			limiter.mu.Lock()
			limiter.ips = make(map[string]*rate.Limiter)
			limiter.mu.Unlock()
		}
	}
}