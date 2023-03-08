package middleware

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt"
)

// Define global variables for email and ID
var email string
var id string

// AdminAuth middleware function for handling admin authentication
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT from cookie
		tokenString, err := c.Cookie("Adminjwt")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}
		// Validate JWT
		err = ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserAuth middleware function for handling user authentication
func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT from cookie
		tokenString, err := c.Cookie("UserAuth")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}
		// Validate JWT
		err = ValidateToken(tokenString)
		// Set email and ID values from JWT claims
		c.Set("user", email)
		c.Set("id", id)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Define JWT key
var JwtKey = []byte("supersecretkey")

// Define JWTClaim struct for storing JWT claims
type JWTClaim struct {
	Email string `json:"email"`
	Uid   uint   `json:"uid"`
	jwt.StandardClaims
}

// GenerateJWT function for generating a new JWT
func GenerateJWT(email string, uid uint) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		Uid:   uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JwtKey)
	return
}

// ValidateToken function for validating a JWT
func ValidateToken(signedToken string) (err error) {
	// Parse JWT with claims
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
	)
	if err != nil {
		return
	}
	// Set email and ID values from JWT claims
	claims, ok := token.Claims.(*JWTClaim)
	id = strconv.Itoa(int(claims.Uid))
	email = claims.Email
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
