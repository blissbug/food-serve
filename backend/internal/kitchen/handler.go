package kitchen

import (
	"fmt"
	"strconv"

	kitchenmembers "food-serve.com/internal/kitchen-members"
	"food-serve.com/internal/middleware"
	"food-serve.com/internal/user"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/validator"
	"github.com/gin-gonic/gin"
)

type KitchenHandler struct {
	KitchenStore       KitchenStore
	KitchenMemberStore kitchenmembers.KitchenMembersStore
	UserStore          *user.UserStore
}

func NewHandler(KitchenStore KitchenStore, KitchenMemberStore kitchenmembers.KitchenMembersStore, UserStore *user.UserStore) *KitchenHandler {
	return &KitchenHandler{
		KitchenStore:       KitchenStore,
		KitchenMemberStore: KitchenMemberStore,
		UserStore:          UserStore,
	}
}

func (h KitchenHandler) CreateKitchen(ctx *gin.Context) {
	var kitchen types.Kitchen
	err := ctx.ShouldBindJSON(&kitchen)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = validator.Validates(kitchen)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	kitchenId, err := h.KitchenStore.CreateKitchen(kitchen.Name, kitchen.Slug)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := ctx.GetUint(middleware.UserIDKey)

	if userId == 0 {
		ctx.JSON(400, gin.H{"error": "User ID not found, please login"})
		return
	}

	err = h.KitchenMemberStore.CreateKitchenMember(userId, kitchenId, "admin")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Kitchen created successfully"})
}

func (h KitchenHandler) SubscribeToKitchenAsUser(ctx *gin.Context) {
	var subscribeToKitchenPayload types.SubscribeKitchenPayload
	err := ctx.ShouldBindJSON(&subscribeToKitchenPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := ctx.GetUint(middleware.UserIDKey)
	fmt.Println(userId)

	if userId == 0 {
		ctx.JSON(400, gin.H{"error": "User ID not found, please login"})
		return
	}

	err = h.KitchenMemberStore.CreateKitchenMember(userId, subscribeToKitchenPayload.KitchenID, "user")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//revoke the token and ask user to login again to gain access to kitchen
	ctx.JSON(200, gin.H{"message": "User subscribed to kitchen successfully"})
}

func (h KitchenHandler) SubscribeToKitchenAsAdmin(ctx *gin.Context) {
	id := ctx.Param("kitchenId")
	fmt.Println(id)

	var subscribeToKitchenAsAdminPayload types.SubscribeKitchenAsAdminPayload
	err := ctx.ShouldBindJSON(&subscribeToKitchenAsAdminPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	kitchenId, err := strconv.Atoi(id)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	value, exists := ctx.Get(middleware.KitchenAndRoleKey)

	if !exists {
		ctx.JSON(401, gin.H{"error": "No kitchen data found!"})
		return
	}

	kitchenMembers, ok := value.([]types.KitchenMemberClaims)

	if !ok {
		ctx.JSON(401, gin.H{"error": "Kitchen data is not in correct format!"})
		return
	}

	var isAdminOfKitchen bool
	for i := 0; i < len(kitchenMembers); i++ {
		kitchenMember := kitchenMembers[i]
		if kitchenMember.KitchenID == uint(kitchenId) && kitchenMember.Role == "admin" {
			isAdminOfKitchen = true
			break
		}
	}

	if !isAdminOfKitchen {
		ctx.JSON(401, gin.H{"error": "User is not admin of this kitchen"})
		return
	}

	u, err := h.UserStore.FindUserByEmail(subscribeToKitchenAsAdminPayload.Email)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.KitchenMemberStore.CreateKitchenMember(u.ID, uint(kitchenId), "admin")

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "User subscribed to kitchen as admin successfully"})

}
