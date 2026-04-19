package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Maker struct {
	secretKey string
}
type Claims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// constructor
func NewMaker(s string) (*Maker, error) {
    if len(s) < 32 {
        return nil, fmt.Errorf("secret key must be at least 32 characters")
    }
    return &Maker{secretKey: s}, nil
}


// generate access token
func (m *Maker) GenerateToken(userId, email string , duration time.Duration) (string, error) {
	claims := Claims{
		UserId: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "second-brain-auth",
			Subject: userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}



// validate the token
func (m *Maker) ValidateToken(tokenString string) (*Claims, error) {

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
