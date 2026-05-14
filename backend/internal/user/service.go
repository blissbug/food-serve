package user

import (
	"fmt"
	"net/http"

	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func (userStore *UserStore) registerUser(payload types.RegisterPayload, ctx *gin.Context) {
	if payload.Email == "" || payload.Password == "" || payload.Username == "" {
		response.Error(ctx, fmt.Errorf("could not parse the request, a required field is empty, in register"), http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
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
		zap.L().Error("could not parse the request, validation error", zap.Error(err))
		response.Error(ctx, fmt.Errorf("could not parse the request, please ensure all required fields are present and username is more than 3 letters"), 400)
		return
	}

	//interface for functions user store ka db has
	err = userStore.CreateUser(user)
	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
		return
	}

	response.JSON(ctx, gin.H{"message": "User created successfully"})
}
