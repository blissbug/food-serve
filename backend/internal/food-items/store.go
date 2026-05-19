package food_items

import (
	"food-serve.com/internal/db"
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

func (fStore FoodItemsStore) UpdateFoodItem(tx *gorm.DB, foodItem map[string]interface{}, itemId uint, kitchenId uint) (error, int64) {
	res := tx.Model(&types.FoodItem{}).Where("id = ? AND kitchen_id = ?", itemId, kitchenId).Updates(foodItem)
	if res.Error != nil {
		return res.Error, 0
	}
	return nil, res.RowsAffected
}

func (fStore FoodItemsStore) GetAllFoodItemsByKitchenId(kitchenId uint, page string, pageSize string) ([]types.FoodItem, error) {
	var foodItems []types.FoodItem
	err := fStore.FoodItemsStore.Scopes(db.Paginate(page, pageSize)).Where("kitchen_id = ?", kitchenId).Find(&foodItems).Error
	if err != nil {
		return nil, err
	}
	return foodItems, nil
}
