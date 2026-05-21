package menu

import (
	"errors"
	"time"

	"food-serve.com/pkg/types"
	"food-serve.com/pkg/validator"
	"gorm.io/gorm"
)

type Service struct {
	MenuStore *Store
}

func NewService(MenuStore *Store) *Service {
	return &Service{MenuStore: MenuStore}
}

func (menuService *Service) CreateMenuService(CreateMenuPayload types.CreateMenuPayload, CreatedBy uint, KitchenId uint) error {
	now := time.Now()

	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		now.Location(),
	)

	menuDate := time.Date(
		CreateMenuPayload.Date.Year(),
		CreateMenuPayload.Date.Month(),
		CreateMenuPayload.Date.Day(),
		0, 0, 0, 0,
		CreateMenuPayload.Date.Location(),
	)

	if menuDate.Before(today) {
		return errors.New("menu date cannot be before today")
	}

	//check status
	if CreateMenuPayload.Status == types.MenuStatusActive {
		//we need the timings then
		//default for now is +7 hours from the current time
		//TODO: clean up this plus make it more readable
		if CreateMenuPayload.OrderOpen == nil {
			now := time.Now()
			closeTime := now.Add(7 * time.Hour)

			CreateMenuPayload.OrderOpen = &now
			CreateMenuPayload.OrderClose = &closeTime
		}
		if CreateMenuPayload.OrderClose == nil && CreateMenuPayload.OrderOpen != nil {
			orderClose := (*CreateMenuPayload.OrderOpen).Add(7 * time.Hour)
			CreateMenuPayload.OrderClose = &orderClose
		}
		if CreateMenuPayload.OrderClose != nil && CreateMenuPayload.OrderOpen == nil {
			return errors.New("please set up order open timings, only order closing times were received")
		}
	}

	if CreateMenuPayload.OrderOpen != nil && CreateMenuPayload.OrderClose != nil {
		if CreateMenuPayload.OrderClose.Before(*CreateMenuPayload.OrderOpen) {
			return errors.New("order close cannot be before order open")
		}
	}

	var FoodItemsList []uint

	for _, i := range CreateMenuPayload.Items {
		FoodItemsList = append(FoodItemsList, uint(i))
	}

	Menu := types.Menu{
		Date:       CreateMenuPayload.Date,
		Status:     CreateMenuPayload.Status,
		OrderOpen:  CreateMenuPayload.OrderOpen,
		OrderClose: CreateMenuPayload.OrderClose,
		UpdatedBy:  CreatedBy,
		CreatedBy:  CreatedBy,
		KitchenID:  KitchenId,
	}

	err := menuService.MenuStore.store.Transaction(func(tx *gorm.DB) error {
		resp := tx.Create(&Menu)
		if resp.Error != nil {
			return resp.Error
		}
		for _, fd := range FoodItemsList {
			var FoodItem types.FoodItem

			resp := tx.Find(&FoodItem, fd)

			if resp.Error != nil {
				return resp.Error
			}

			MenuItem := types.MenuItem{
				MenuID:      Menu.ID,
				FoodItemID:  fd,
				Price:       FoodItem.Price,
				IsAvailable: true,
			}
			resp = tx.Create(&MenuItem)
			if resp.Error != nil {
				return resp.Error
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// PublishMenuService runs only on the first publish action that you edit after
func (menuService *Service) PublishMenuService(PublishMenuPayload types.PublishMenuPayload, menuId uint, kitchenId uint) error {
	//we need the timings then
	//default for now is +7 hours from the current time

	validationErr := validator.Validates(PublishMenuPayload)

	if validationErr != nil {
		return validationErr
	}

	menu, exists, err := menuService.MenuStore.GetMenu(menuId, kitchenId)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no such menu exists in your kitchen")
	}

	if menu.Status == types.MenuStatusActive {
		return errors.New("menu is already active")
	}

	if menu.OrderOpen == nil {
		return errors.New("order open timings are required")
	}
	if menu.OrderClose == nil {
		return errors.New("order close timings are required")
	}

	dbErr := menuService.MenuStore.PublishMenu(menuId)

	if dbErr != nil {
		return dbErr
	}

	return nil
}

func (menuService *Service) UpdateMenuService(UpdateMenuPayload types.UpdateMenuPayload, menuId uint, kitchenId uint) error {
	menu, exists, err := menuService.MenuStore.GetMenu(menuId, kitchenId)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no such menu exists in your kitchen")
	}

	isDraft := menu.Status == types.MenuStatusDraft

	var UpdatedMenu types.Menu

	if UpdateMenuPayload.Date != nil {
		if !isDraft {
			return errors.New("cannot update date of published menu")
		}
		UpdatedMenu.Date = *UpdateMenuPayload.Date
	}

	if UpdateMenuPayload.OrderOpen != nil {
		if !isDraft {
			return errors.New("cannot update order open timings of published menu")
		}
		UpdatedMenu.OrderOpen = UpdateMenuPayload.OrderOpen
	}

	if UpdateMenuPayload.OrderClose != nil {
		if !isDraft {
			if UpdateMenuPayload.OrderClose.Before(time.Now()) {
				return errors.New("cannot update order close timings to a time in the past")
			}
		}
		UpdatedMenu.OrderClose = UpdateMenuPayload.OrderClose
	}

	transactionErr := menuService.MenuStore.store.Transaction(func(tx *gorm.DB) error {
		var menuItems []types.MenuItem
		resp := tx.Where("menu_id = ?", menuId).Find(&menuItems)
		if resp.Error != nil {
			return resp.Error
		}

		existingItems := make(map[uint]types.MenuItem, len(menuItems))
		for _, item := range menuItems {
			existingItems[item.FoodItemID] = item
		}
		for _, item := range UpdateMenuPayload.Items {
			var foodItem types.FoodItem
			//TODO: improve this to one call to db
			foodItemResp := tx.
				Where(
					"id=? AND kitchen_id=?",
					item,
					kitchenId,
				).
				First(&foodItem)

			if foodItemResp.Error != nil {
				return foodItemResp.Error
			}

			if existingItem, ok := existingItems[uint(item)]; ok {
				if isDraft {
					//update with price
					resp := tx.
						Model(&types.MenuItem{}).
						Where(
							"menu_id=? AND food_item_id=?",
							menuId,
							item,
						).
						Update("price", foodItem.Price)
					if resp.Error != nil {
						return resp.Error
					}
				}
				if !isDraft && !existingItem.IsAvailable {

					if !foodItem.IsActive {
						return errors.New("cannot reactivate inactive item")
					}

					resp := tx.
						Model(&types.MenuItem{}).
						Where(
							"menu_id=? AND food_item_id=?",
							menuId,
							item,
						).
						Update("is_available", true)

					if resp.Error != nil {
						return resp.Error
					}
				}
				delete(existingItems, uint(item))
			} else {

				if !foodItem.IsActive {
					return errors.New("cannot add inactive food item to menu")
				}

				menuItem := types.MenuItem{
					MenuID:      menuId,
					FoodItemID:  uint(item),
					IsAvailable: true,
					Price:       foodItem.Price,
				}
				resp := tx.Create(&menuItem)

				if resp.Error != nil {
					return resp.Error
				}
			}
		}

		for _, item := range existingItems {
			if isDraft {
				resp := tx.Delete(&item)
				if resp.Error != nil {
					return resp.Error
				}
			} else {
				resp := tx.Model(&item).Update("is_available", false)
				if resp.Error != nil {
					return resp.Error
				}
			}
		}
		resp = tx.Model(&menu).Updates(&UpdatedMenu)
		if resp.Error != nil {
			return resp.Error
		}
		return nil
	})

	if transactionErr != nil {
		return transactionErr
	}

	return nil
}

func (menuService *Service) DeleteMenuService(menuId uint, kitchenId uint) error {
	menu, exists, err := menuService.MenuStore.GetMenu(menuId, kitchenId)

	if err != nil {
		return err
	}

	if !exists {
		return errors.New("no such menu exists in your kitchen")
	}

	if menu.Status == types.MenuStatusActive {
		return errors.New("menu is active, cannot delete")
	}

	dbErr := menuService.MenuStore.DeleteMenu(menuId, kitchenId)

	if dbErr != nil {
		return dbErr
	}

	return nil
}
