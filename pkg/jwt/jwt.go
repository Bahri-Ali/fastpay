package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(UserID string , Role string , secret string , expirationHours int)(string , error){

	Claims:= Claims{
		UserID: UserID,
		Role: Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationHours)*time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "FASTPAY",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 ,Claims )
	return  token.SignedString([]byte(secret))
}

// func ValidateToken(Token string , secret string ) (*Claims , error){
// 	token , err := jwt.ParseWithClaims(Token , &Claims{} , func(token *jwt.Token)(interface{}, error){
// 		return []byte(secret), nil
// 	})
// 	if err != nil {return  nil,err}

// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//         return claims, nil
//     }
// 	return   nil, fmt.Errorf("invalid token")
// }
func ValidateToken() string {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}