package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaim struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(email, id, role, secret string) (string, error) {
	now := time.Now()
	claims := &JWTClaim{
		Email: email,
		ID:    id,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(signedToken, secret string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

func AdminAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("Adminjwt")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
			c.Abort()
			return
		}
		claims, err := ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"status": false, "message": "Admin access required"})
			c.Abort()
			return
		}
		c.Set("admin_email", claims.Email)
		c.Set("admin_id", claims.ID)
		c.Next()
	}
}

func UserAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("UserAuth")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Unauthorized"})
			c.Abort()
			return
		}
		claims, err := ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"status": false, "message": "User access required"})
			c.Abort()
			return
		}
		c.Set("user_email", claims.Email)
		c.Set("user_id", claims.ID)
		c.Next()
	}
}
