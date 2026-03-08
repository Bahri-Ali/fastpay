package auth

type RegisterRequest struct {
    PhoneNumber string `json:"phone_number" binding:"required"`
    Password    string `json:"password" binding:"required,min=6"`
    FullName    string `json:"full_name" binding:"required"`
    Email       string `json:"email"` 
    
    WilayaCode int `json:"wilaya_code" binding:"required"`
    

}

type LoginRequest struct {
    PhoneNumber string `json:"phone_number" binding:"required"`
    Password    string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token         string `json:"token"`
}