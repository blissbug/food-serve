package user

import (
	"fmt"

	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Handler depending on requirement I can have multiple dbs of multiple stores as userStore,
// FoodStore all are of same type and can be extended accordingly
// this allows to inject multiple dbs
type Handler struct {
	userStore *UserStore //has db internally
}

func NewHandler(userStore *UserStore) *Handler {
	return &Handler{
		userStore: userStore,
	}
}

func (h Handler) HandleLogin(ctx *gin.Context) {
	var loginPayload types.LoginPayload
	err := ctx.ShouldBindJSON(&loginPayload) //bind body from request to the object

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginPayload)

	if err != nil {
		response.Error(ctx, fmt.Errorf("could not parse the request, a required field is empty"), 400)
		return
	}

	user, err := h.userStore.FindUser(loginPayload.Email)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	isPasswordValid := utils.ComparePassword(loginPayload.Password, user.Password)

	if !isPasswordValid {
		response.Error(ctx, fmt.Errorf("invalid credentials"), 400)
		return
	}

	response.JSON(ctx, gin.H{"message": "Login successful"})
}

func (h Handler) HandleRegister(ctx *gin.Context) {
	//parse request and get username, email and password
	var registerPayload types.RegisterPayload
	err := ctx.ShouldBindJSON(&registerPayload)

	if err != nil {
		response.Error(ctx, err, 400)
		return
	}

	h.userStore.registerUser(registerPayload, ctx)
}
