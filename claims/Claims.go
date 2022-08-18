package claims

import "github.com/golang-jwt/jwt"

type MapClaims struct {
	SessionID string `json:"sessionId"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	jwt.StandardClaims
}
