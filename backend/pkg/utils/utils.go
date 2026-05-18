package utils

import (
	"fmt"
	"time"

	"food-serve.com/pkg/config"
	"food-serve.com/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPasswordInBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedPasswordInBytes)

	return hashedPassword, nil
}

func ComparePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func JWTwithClaimsForTokens(user types.User, kitchenMember []types.KitchenMember) (string, string, error) {
	Env, err := config.LoadConfig()
	fmt.Println(Env)

	if err != nil {
		return "", "", err
	}

	var accessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":        user.ID,
		"email":         user.Email,
		"kitchenMember": kitchenMember,
		"exp":           time.Now().Add(time.Minute * 15).Unix(),
		"iat":           time.Now().Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(Env.JWTSecret))

	var refreshToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), //goes for a week
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(Env.JWTSecret))

	return accessTokenString, refreshTokenString, err
}
