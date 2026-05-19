package food_items

import (
	"fmt"
	"os"
	"path/filepath"

	"food-serve.com/internal/middleware"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"food-serve.com/pkg/validator"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageDataStruct struct {
	ImageURL            string `json:"imageUrl"`
	OriginalDestination string `json:"originalDestination"`
	OriginalOrder       int    `json:"originalOrder"`
}

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

	var ItemImageData []ImageDataStruct

	for order, file := range files {
		err := os.MkdirAll("./assets/uploads", os.ModePerm)
		if err != nil {
			return
		}
		// Upload the file to specific dst.
		dst := filepath.Join("./assets/uploads", fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(file.Filename)))
		fileSaveError := ctx.SaveUploadedFile(file, dst)

		if fileSaveError != nil {
			ctx.JSON(400, gin.H{"error": fileSaveError.Error()})
			return
		}

		//cloudinary job
		resp, uploadErr := h.cld.Upload.Upload(ctx, dst, uploader.UploadParams{PublicID: uuid.NewString()})

		if uploadErr != nil {
			ctx.JSON(400, gin.H{"error": uploadErr.Error()})
			return
		}
		//save with old file name
		fmt.Println(resp.URL)

		ItemImageData = append(ItemImageData, ImageDataStruct{
			OriginalDestination: dst,
			ImageURL:            resp.URL,
			OriginalOrder:       order,
		})

		//delete it from local storage
		fileRemoveErr := os.Remove(dst)
		if fileRemoveErr != nil {
			ctx.JSON(400, gin.H{"error": fileRemoveErr.Error()})
			return
		}
	}

	creationError := h.FoodItemsService.CreateFoodItem(foodItem, ItemImageData)

	if creationError != nil {
		ctx.JSON(400, gin.H{"error": creationError.Error(), "message": "Food item creation failed, please try again"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Food item created successfully"})
}
