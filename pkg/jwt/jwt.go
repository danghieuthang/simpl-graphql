package jwt

import (
	"example/web-service-gin/entity"
	"log"
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
	claims["name"] = user.Name
	claims["sub"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// ParseToken parses a jwt token and returns the user in it's claims
func ParseToken(tokenStr string) (*entity.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		name := claims["name"].(string)
		id := int(claims["sub"].(float64))
		return &entity.User{
			Email: email,
			Name:  name,
			Id:    id,
		}, nil
	} else {
		return nil, err
	}
}
