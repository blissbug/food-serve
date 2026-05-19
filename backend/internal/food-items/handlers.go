package food_items

import (
	"fmt"
	"strconv"

	"food-serve.com/internal/middleware"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"food-serve.com/pkg/validator"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	FoodItemsService *FoodService
	cld              *cloudinary.Cloudinary
}

func NewHandler(FoodItemService *FoodService, cld *cloudinary.Cloudinary) *Handler {
	return &Handler{
		FoodItemsService: FoodItemService,
		cld:              cld,
	}
}

func (h Handler) CreateFoodItem(ctx *gin.Context) {
	var CreateFoodItemPayload types.CreateFoodItemPayload
	err := ctx.ShouldBind(&CreateFoodItemPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = validator.Validates(CreateFoodItemPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get kitchen members for the said kitchen
	value, exists := ctx.Get(middleware.KitchenAndRoleKey)

	if !exists {
		ctx.JSON(401, gin.H{"error": fmt.Errorf("user does not have access to this kitchen as admin")})
		return
	}

	fmt.Println(value)

	kitchensAndRoles, ok := value.([]types.KitchenAndRoleClaims)

	if !ok {
		fmt.Println("kitchen and roles not in correct format")
	}

	//check if the user is admin of the said kitchen
	isUserAdminOfThisKitchen := utils.CheckRoleForThisKitchen(CreateFoodItemPayload.KitchenID, "admin", kitchensAndRoles)

	if !isUserAdminOfThisKitchen {
		ctx.JSON(401, gin.H{"error": fmt.Errorf("user does not have access to this kitchen as admin")})
		return
	}
	//get the files
	form, err := ctx.MultipartForm()

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	foodItem := types.FoodItem{
		Name:        CreateFoodItemPayload.Name,
		Description: CreateFoodItemPayload.Description,
		Price:       CreateFoodItemPayload.Price,
		KitchenID:   CreateFoodItemPayload.KitchenID,
		IsActive:    CreateFoodItemPayload.IsActive,
		IsVeg:       CreateFoodItemPayload.IsVeg,
		IsVegan:     CreateFoodItemPayload.IsVegan,
	}

	files := form.File["images"]

	ItemImageData, uploadErr := utils.UploadImagesToCloudinary(files, h.cld, ctx)

	if uploadErr != nil {
		ctx.JSON(400, gin.H{"error": uploadErr.Error()})
		return
	}

	creationError := h.FoodItemsService.CreateFoodItem(foodItem, ItemImageData)

	if creationError != nil {
		ctx.JSON(400, gin.H{"error": creationError.Error(), "message": "Food item creation failed, please try again"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Food item created successfully"})
}

func (h Handler) UpdateFoodItem(ctx *gin.Context) {
	var UpdateFoodItemPayload types.UpdateFoodItemPayload
	err := ctx.ShouldBind(&UpdateFoodItemPayload)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = validator.Validates(UpdateFoodItemPayload)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	value := ctx.Param("foodItemId")
	foodItemId, err := strconv.Atoi(value)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//get the files
	form, err := ctx.MultipartForm()

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	files := form.File["images"]

	ItemImageData, uploadErr := utils.UploadImagesToCloudinary(files, h.cld, ctx)

	if uploadErr != nil {
		ctx.JSON(400, gin.H{"error": uploadErr.Error()})
		return
	}

	var existingImageIds []uint

	for _, value := range UpdateFoodItemPayload.ExistingImageIDs {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		existingImageIds = append(existingImageIds, uint(intValue))
	}

	err = h.FoodItemsService.UpdateFoodItem(UpdateFoodItemPayload, uint(foodItemId), ItemImageData, existingImageIds)

	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "Food item updated successfully"})
}
