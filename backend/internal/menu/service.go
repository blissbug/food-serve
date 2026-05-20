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
