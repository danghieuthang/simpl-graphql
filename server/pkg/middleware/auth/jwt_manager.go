package auth

import (
	"example/web-service-gin/pkg/entity"
	"example/web-service-gin/pkg/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// secret key being used to sign tokens
var SecretKey = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken generates a jwt token and assign a user information to it's claims and return it
func GenerateToken(user *entity.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["email"] = user.Email
	claims["role"] = user.Role.Name
	claims["name"] = user.Name
	claims["sub"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		logger.Logger.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the user in it's claims
func ParseToken(tokenStr string) (*AuthenticatedUser, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		role := claims["role"].(string)
		name := claims["name"].(string)
		id := int(claims["sub"].(float64))
		return &AuthenticatedUser{
			Email: email,
			Name:  name,
			Id:    id,
			Role:  role,
		}, nil
	} else {
		return nil, err
	}
}
