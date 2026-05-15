package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"food-serve.com/internal/cache"
	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"food-serve.com/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	otpTypeLogin    = "login"
	otpTypeRegister = "register"
)

func (h Handler) registerUser(payload types.RegisterPayload, ctx *gin.Context) {
	//validate email and required fields as present in the struct
	err := validator.Validates(payload)

	if err != nil {
		zap.L().Error("could not parse the request, validation error", zap.Error(err))
		response.Error(ctx, err, 400)
		return
	}

	existingUser, err := h.userStore.FindUserByEmail(payload.Email)

	if err == nil && existingUser.IsVerified {
		response.Error(ctx, fmt.Errorf("user with this email already exists, please login"), http.StatusBadRequest)
		return
	}

	if err == nil && existingUser.IsVerified == false {
		id := existingUser.ID
		err = generateOTPAndSetToRedis(ctx, h.rdb, id, otpTypeRegister)
		if err != nil {
			response.Error(ctx, err, http.StatusBadRequest)
			return
		}
		response.Error(ctx, fmt.Errorf("user with this email already exists, resent the otp, please verify your email"), http.StatusBadRequest)
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

	//interface for functions user store ka db has
	id, err := h.userStore.CreateUser(user)
	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
		return
	}

	err = generateOTPAndSetToRedis(ctx, h.rdb, id, otpTypeRegister)
	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
		return
	}

	response.JSON(ctx, gin.H{"message": "User created successfully, otp sent to your email, please verify"})
}

func generateOTP() (string, string, error) {
	val := rand.Intn(1000000)
	otp := fmt.Sprintf("%06d", val)
	//send otp in mail to user
	hashedOtp, err := bcrypt.GenerateFromPassword([]byte(otp), 12)

	return otp, string(hashedOtp), err
}

func generateOTPAndSetToRedis(ctx *gin.Context, rdb cache.Cache, id uint, otpType string) error {
	redisKey := fmt.Sprintf("%d:otp:%v", id, otpType)
	fmt.Println(redisKey)
	otp, hashedOtp, err := generateOTP()

	fmt.Println(otp) //we send this to user

	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
		return err
	}
	//store otp in redis
	err = rdb.SetKey(ctx, redisKey, hashedOtp, time.Minute*10)

	if err != nil {
		response.Error(ctx, err, http.StatusBadRequest)
		return err
	}
	return nil
}
