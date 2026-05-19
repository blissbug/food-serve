package middleware

import (
	"errors"
	"fmt"
	"strings"

	"food-serve.com/pkg/config"
	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID           uint `json:"userId"`
	Email            string
	Role             string
	KitchensAndRoles []types.KitchenAndRoleClaims `json:"kitchenMember"`
	jwt.RegisteredClaims
}

const (
	UserIDKey         = "userID"
	EmailKey          = "email"
	RoleKey           = "role"
	KitchenAndRoleKey = "kitchenAndRole"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Env, err := config.LoadConfig()
		if err != nil {
			response.Error(ctx, err, 500)
			return
		}

		tokenString := ctx.GetHeader("Authorization")
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(ctx, errors.New("invalid token"), 401)
			ctx.Abort()
			return
		}

		tokenString = parts[1]

		fmt.Println("i was here??")

		if tokenString == "" {
			response.Error(ctx, errors.New("unauthorized"), 401)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(Env.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			response.Error(ctx, err, 401)
			ctx.Abort()
			return
		} else if claims, ok := token.Claims.(*JWTClaims); ok {
			ctx.Set(UserIDKey, claims.UserID)
			ctx.Set(EmailKey, claims.Email)
			ctx.Set(RoleKey, claims.Role)
			ctx.Set(KitchenAndRoleKey, claims.KitchensAndRoles)
			fmt.Println(claims.KitchensAndRoles)
			fmt.Println(claims.UserID, claims.Email, claims.Role)
			ctx.Next()
		} else {
			response.Error(ctx, errors.New("invalid token"), 401)
			ctx.Abort()
			return
		}
	}
}
