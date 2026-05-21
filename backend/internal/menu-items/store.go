package menuItems

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type Store struct {
	Store *gorm.DB
}

func NewStore(store *gorm.DB) *Store {
	return &Store{Store: store}
}

func (MenuItemsStore *Store) GetMenuItemsForMenuId(menuId uint) ([]types.MenuItem, error) {
	var menuItems []types.MenuItem
	resp := MenuItemsStore.Store.Where("menu_id = ?", menuId).Find(&menuItems)

	if resp.Error != nil {
		return menuItems, resp.Error
	}

	return menuItems, nil
}
