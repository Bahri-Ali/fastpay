package user

type ProfileResponse struct {
    ID             string `json:"id"`
    UserIdentifier string `json:"user_identifier"`
    PhoneNumber    string `json:"phone_number"`
    Email          string `json:"email"`
    FullName       string `json:"full_name"`
    Role           string `json:"role"`
}

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=6"`
}

type VerifyPasswordRequest struct {
    VerificationToken string `json:"verification_token" binding:"required"`
    PIN               string `json:"pin" binding:"required"`
}

type ActionResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    VerificationToken string `json:"verification_token,omitempty"` 
}