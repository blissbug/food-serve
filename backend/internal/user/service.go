package user

import (
	"net/http"

	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (userStore *UserStore) registerUser(payload types.RegisterPayload, ctx *gin.Context) {
	if payload.Email == "" || payload.Password == "" || payload.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request, a required field is empty"})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request"})
		return
	}

	user := types.User{
		Email:    payload.Email,
		Password: hashedPassword,
		Username: payload.Username,
		Role:     payload.Role,
	}

	//validate email and required fields as present in the struct
	validate := validator.New()
	err = validate.Struct(user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse the request, a required field is empty"})
		return
	}

	//interface for functions user store ka db has
	err = userStore.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
