package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func SignToken(userId string, username, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresIn := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userId,
		"user": username,
		"role": role,
	}

	if jwtExpiresIn != "" {
		duration, err := time.ParseDuration(jwtExpiresIn)
		if err != nil {
			return "", ErrorHandler(err, "Internal error")
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", ErrorHandler(err, "Internal error")
	}

	return signedToken, nil
}
