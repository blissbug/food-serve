package food_items

import (
	foodItemImages "food-serve.com/internal/food-item-images"
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type FoodService struct {
	db         *gorm.DB
	foodStore  FoodItemsStore
	imageStore foodItemImages.Store
}

func NewService(db *gorm.DB, foodStore FoodItemsStore, imageStore foodItemImages.Store) *FoodService {
	return &FoodService{
		db:         db,
		foodStore:  foodStore,
		imageStore: imageStore,
	}
}

func (f FoodService) CreateFoodItem(foodItem types.FoodItem, ItemImageData []ImageDataStruct) error {
	return f.db.Transaction(func(tx *gorm.DB) error {
		foodItemId, err := f.foodStore.CreateFoodItem(tx, foodItem)

		if err != nil {
			return err
		}

		for _, imageData := range ItemImageData {
			foodItemImage := types.FoodItemImage{
				FoodItemID:          foodItemId,
				ImageURL:            imageData.ImageURL,
				DisplayOrder:        imageData.OriginalOrder,
				OriginalDestination: imageData.OriginalDestination,
			}
			err = f.imageStore.AddImageForItem(tx, foodItemImage)
			if err != nil {
				return err
			}
		}
		return nil
	})

}
