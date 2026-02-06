package utils

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TransactionHistory struct {
	UserId        string
	Amount        float32
	Description   string
	PaymentMethod string
}
type JwtSessionPayload struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

func ValidateToken(tokenString string) (JwtSessionPayload, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	var jwtKey = []byte(jwtSecret)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return JwtSessionPayload{}, err
	}

	claims = token.Claims.(jwt.MapClaims)

	jsonString, _ := json.Marshal(claims)

	jwtPayload := JwtSessionPayload{}
	json.Unmarshal(jsonString, &jwtPayload)

	return jwtPayload, err
}

func SignToken(user_id string, email string) (string, error) {

	token_lifespan := time.Now().Add(time.Hour * 24).Unix()

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["email"] = email
	claims["exp"] = token_lifespan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GetUserIdFromToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no metadata found")
	}

	val, ok := md["authorization"]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	tokenString := strings.TrimPrefix(val[0], "Bearer ")
	if tokenString == "" {
		return "", status.Error(codes.Unauthenticated, "Unauthorized Access")
	}

	return tokenString, nil
}
