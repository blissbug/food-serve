package user

import (
	"fmt"

	"food-serve.com/internal/cache"
	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"food-serve.com/pkg/validator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Handler depending on requirement, I can have multiple dbs of multiple stores as userStore,
// FoodStore all are of the same type and can be extended accordingly
// this allows injecting multiple dbs
type Handler struct {
	userStore *UserStore //has db internally
	rdb       cache.Cache
}

func NewHandler(userStore *UserStore, rdb cache.Cache) *Handler {
	return &Handler{
		userStore: userStore,
		rdb:       rdb,
	}
}

func (h Handler) HandleLogin(ctx *gin.Context) {
	var loginPayload types.LoginPayload
	err := ctx.ShouldBindJSON(&loginPayload) //bind body from request to the object

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	err = validator.Validates(loginPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	user, err := h.userStore.FindUserByEmail(loginPayload.Email)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	isPasswordValid := utils.ComparePassword(loginPayload.Password, user.Password)

	if !isPasswordValid {
		response.Error(ctx, fmt.Errorf("invalid credentials"), 400)
		return
	}

	err = generateOTPAndSetToRedis(ctx, h.rdb, user.ID, otpTypeLogin)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	response.JSON(ctx, gin.H{"message": "User logged in successfully, otp sent to your email, please verify"})
}

func (h Handler) HandleRegister(ctx *gin.Context) {
	//parse request and get username, email and password
	var registerPayload types.RegisterPayload
	err := ctx.ShouldBindJSON(&registerPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	h.registerUser(registerPayload, ctx)
}

func (h Handler) HandleVerifyRegistration(ctx *gin.Context) {
	//check if otp is valid
	var RegistrationOTPPayload types.OTP
	err := ctx.ShouldBindJSON(&RegistrationOTPPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	fmt.Println(RegistrationOTPPayload)

	err = validator.Validates(RegistrationOTPPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	user, err := h.userStore.FindUserByEmail(RegistrationOTPPayload.Email)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	if user.ID == 0 {
		response.Error(ctx, fmt.Errorf("user not found"), 400)
		return
	}

	redisKey := fmt.Sprintf("%d:otp:%v", user.ID, otpTypeRegister)
	fmt.Println(redisKey)
	hashedOtp, err := h.rdb.GetKey(ctx, redisKey)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	valid := bcrypt.CompareHashAndPassword([]byte(hashedOtp), []byte(RegistrationOTPPayload.Otp))

	if valid != nil {
		response.Error(ctx, fmt.Errorf("invalid otp"), 400)
		return
	}
	//validated
	//if valid make is verified to true by updating the user
	validated, err := h.userStore.VerifyUser(user)
	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	if !validated {
		response.Error(ctx, fmt.Errorf("invalid otp"), 400)
	}

	response.JSON(ctx, gin.H{"message": "User verified successfully"})
}

func (h Handler) HandlerVerifyLogin(ctx *gin.Context) {
	var loginOTPPayload types.OTP
	err := ctx.ShouldBindJSON(&loginOTPPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}
	//from otp check if its valid

	user, err := h.userStore.FindUserByEmail(loginOTPPayload.Email)
	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	redisKey := fmt.Sprintf("%d:otp:%v", user.ID, otpTypeLogin)
	hashedOtp, err := h.rdb.GetKey(ctx, redisKey)
	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	valid := bcrypt.CompareHashAndPassword([]byte(hashedOtp), []byte(loginOTPPayload.Otp))
	if valid != nil {
		response.Error(ctx, fmt.Errorf("invalid otp"), 400)
		return
	}
	token, err := utils.JWTwithClaims(user)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	response.JSON(ctx, gin.H{"message": "Login successful", "token": token})
}
