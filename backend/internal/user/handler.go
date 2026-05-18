package user

import (
	"fmt"
	"time"

	"food-serve.com/internal/cache"
	kitchenmembers "food-serve.com/internal/kitchen-members"
	"food-serve.com/pkg/config"
	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"food-serve.com/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Handler depending on requirement, I can have multiple dbs of multiple stores as userStore,
// FoodStore all are of the same type and can be extended accordingly
// this allows injecting multiple dbs
type Handler struct {
	userStore          *UserStore //has db internally
	kitchenMemberStore *kitchenmembers.KitchenMembersStore
	rdb                cache.Cache
}

type RefreshTokenClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

func NewHandler(userStore *UserStore, rdb cache.Cache, kitchenMemberStore *kitchenmembers.KitchenMembersStore) *Handler {
	return &Handler{
		userStore:          userStore,
		rdb:                rdb,
		kitchenMemberStore: kitchenMemberStore,
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

	kitchenMember, err := h.kitchenMemberStore.GetKitchensSubscribedByAUser(user.ID)

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

	accessToken, refreshToken, err := utils.JWTwithClaimsForTokens(user, kitchenMember)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	redisKey = fmt.Sprintf("%d:refreshToken", user.ID)
	//add refresh token to redis
	err = h.rdb.SetKey(ctx, redisKey, refreshToken, time.Hour*24*7)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	response.JSON(ctx, gin.H{"message": "Login successful", "access-token": accessToken, "refresh-token": refreshToken})
}

func (h Handler) HandleRefresh(ctx *gin.Context) {
	var RefreshPayload types.RefreshToken
	err := ctx.ShouldBindJSON(&RefreshPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	Env, err := config.LoadConfig()

	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	refreshToken, err := jwt.ParseWithClaims(RefreshPayload.RefreshToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(Env.JWTSecret), nil
	})

	if err != nil || !refreshToken.Valid {
		response.Error(ctx, err, 401)
		return
	}

	claims, ok := refreshToken.Claims.(*RefreshTokenClaims)
	if !ok || !refreshToken.Valid {
		response.Error(ctx, err, 401)
		return
	}

	userId := claims.UserID
	redisKey := fmt.Sprintf("%d:refreshToken", userId)

	storedRefreshToken, err := h.rdb.GetKey(ctx, redisKey)

	if err != nil {
		response.Error(ctx, err, 401)
		return
	}

	if storedRefreshToken != RefreshPayload.RefreshToken {
		response.Error(ctx, err, 401)
		return
	}
	//now we have the user id and refresh token which are valid, so we give a new access token + update the refresh token\

	user, err := h.userStore.GetUserById(userId)

	if err != nil {
		response.Error(ctx, err, 401)
		return
	}

	kitchenMember, err := h.kitchenMemberStore.GetKitchensSubscribedByAUser(user.ID)

	if err != nil {
		response.Error(ctx, err, 401)
		return
	}

	newAccessToken, newRefreshToken, err := utils.JWTwithClaimsForTokens(user, kitchenMember)

	if err != nil {
		response.Error(ctx, err, 401)
		return
	}

	err = h.rdb.DeleteKey(ctx, redisKey)
	if err != nil {
		response.Error(ctx, err, 401)
		return
	}
	err = h.rdb.SetKey(ctx, redisKey, newRefreshToken, time.Hour*24*7)

	if err != nil {
		response.Error(ctx, err, 401)
		return
	}

	response.JSON(ctx, gin.H{"message": "Refresh successful", "access-token": newAccessToken, "refresh-token": newRefreshToken})
}
