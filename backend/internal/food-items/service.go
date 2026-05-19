package food_items

import (
	"errors"
	"slices"

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

func (f FoodService) CreateFoodItem(foodItem types.FoodItem, ItemImageData []types.ImageDataStruct) error {
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

func (f FoodService) UpdateFoodItem(updatedFoodItem types.UpdateFoodItemPayload, itemId uint, ItemImageData []types.ImageDataStruct, existingImageIds []uint, kitchenId uint) error {
	var foodItemUpdates = make(map[string]interface{})

	if updatedFoodItem.Name != nil {
		foodItemUpdates["name"] = *updatedFoodItem.Name
	}
	if updatedFoodItem.Description != nil {
		foodItemUpdates["description"] = *updatedFoodItem.Description
	}
	if updatedFoodItem.Price != nil {
		foodItemUpdates["price"] = *updatedFoodItem.Price
	}
	if updatedFoodItem.IsActive != nil {
		foodItemUpdates["is_active"] = *updatedFoodItem.IsActive
	}
	if updatedFoodItem.IsVeg != nil {
		foodItemUpdates["is_veg"] = *updatedFoodItem.IsVeg
	}
	if updatedFoodItem.IsVegan != nil {
		foodItemUpdates["is_vegan"] = *updatedFoodItem.IsVegan
	}

	return f.db.Transaction(func(tx *gorm.DB) error {
		itemImages, err := f.imageStore.GetImagesForItem(f.db, itemId)

		if err != nil {
			return err
		}

		var deletedImages []uint // these images will be deleted
		for _, image := range itemImages {
			if !slices.Contains(existingImageIds, image.ID) {
				deletedImages = append(deletedImages, image.ID)
			}
		}

		err, resp := f.foodStore.UpdateFoodItem(tx, foodItemUpdates, itemId, kitchenId)
		if err != nil {
			return err
		}
		if resp == 0 {
			return errors.New("no such item found for your kitchen")
		}
		//delete the images that are not in the updated payload
		for _, imageId := range deletedImages {
			deleteErr := tx.Delete(&types.FoodItemImage{}, "id = ?", imageId).Error

			if deleteErr != nil {
				return deleteErr
			}
		}

		//add the new images
		for _, imageData := range ItemImageData {
			foodItemImage := types.FoodItemImage{
				FoodItemID:          itemId,
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

func (f FoodService) DeleteFoodItem(itemId uint, kitchenId uint) error {
	return f.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&types.FoodItem{}, "id = ? AND kitchen_id = ?", itemId, kitchenId).Error
		if err != nil {
			return err
		}

		//TODO: delete from cloudinary as well
		return tx.Delete(&types.FoodItemImage{}, "food_item_id = ?", itemId).Error
	})
}

// GetAllFoodItemsForAKitchen gets all the items, active or inactive, to be shown on admin side
func (f FoodService) GetAllFoodItemsForAKitchen(kitchenId uint, page string, pageSize string) ([]types.FoodItem, error) {
	foodItems, err := f.foodStore.GetAllFoodItemsByKitchenId(kitchenId, page, pageSize)

	if err != nil {
		return nil, err
	}
	return foodItems, nil
}
