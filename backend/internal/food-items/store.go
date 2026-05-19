package food_items

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type FoodItemsStore struct {
	FoodItemsStore *gorm.DB
}

func NewStore(store *gorm.DB) FoodItemsStore {
	return FoodItemsStore{
		FoodItemsStore: store,
	}
}

func (fStore FoodItemsStore) CreateFoodItem(tx *gorm.DB, foodItem types.FoodItem) (uint, error) {
	res := tx.Create(&foodItem)
	if res.Error != nil {
		return 0, res.Error
	}
	return foodItem.ID, nil
}
