package user

import (
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// depending on requirement i can have multiple dbs of multiple stores as userStore,
// FoodStore all of same type and can be extended accordingly
// this allows to inject mulitple dbs
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
		ctx.JSON(400, gin.H{"message": "Could not parse the request"})
		return
	}

	validate := validator.New()
	err = validate.Struct(loginPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "Could not parse the request, a required field is empty"})
		return
	}

	user, err := h.userStore.FindUser(loginPayload.Email)

	if err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	isPasswordValid := utils.ComparePassword(loginPayload.Password, user.Password)

	if !isPasswordValid {
		ctx.JSON(400, gin.H{"message": "Invalid credentials"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Login successful"})
}

func (h Handler) HandleRegister(ctx *gin.Context) {
	//parse request and get username, email and password
	var registerPayload types.RegisterPayload
	err := ctx.ShouldBindJSON(&registerPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "Could not parse the request"})
		return
	}

	h.userStore.registerUser(registerPayload, ctx)
}
