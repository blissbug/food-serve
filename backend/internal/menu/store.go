package menu

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type Store struct {
	store *gorm.DB
}

func NewStore(store *gorm.DB) *Store {
	return &Store{store: store}
}

func (menuStore *Store) CreateMenu(MenuItem types.Menu) (uint, error) {
	resp := menuStore.store.Create(&MenuItem)
	if resp.Error != nil {
		return 0, resp.Error
	}

	return MenuItem.ID, nil
}

func (menuStore *Store) GetMenu(menuId uint, kitchenId uint) (types.Menu, bool, error) {
	var menu types.Menu
	resp := menuStore.store.Find(&menu, "id = ? && kitchen_id = ?", menuId, kitchenId)
	if resp.Error != nil {
		return menu, false, resp.Error
	}
	return menu, resp.RowsAffected != 0, nil
}

func (menuStore *Store) UpdateMenu(MenuItem types.Menu) error {
	resp := menuStore.store.Model(&MenuItem).Updates(&MenuItem)
	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (menuStore *Store) PublishMenu(menuId uint) error {
	resp := menuStore.store.Model(&types.Menu{}).Where("id = ?", menuId).Update("status", types.MenuStatusActive)

	if resp.Error != nil {
		return resp.Error
	}
	return nil
}

func (menuStore *Store) DeleteMenu(menuId uint, kitchenId uint) error {
	resp := menuStore.store.Where("id = ? && kitchen_id = ?", menuId, kitchenId).Delete(&types.Menu{})

	if resp.Error != nil {
		return resp.Error
	}
	return nil
}
