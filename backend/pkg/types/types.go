package types

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uint   `json:"id"`
	Username   string `json:"username" validate:"required,min=3,max=20"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

type RegisterPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Role     string `json:"role" validate:"required,oneof=admin user"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SubscribeKitchenPayload struct {
	KitchenID uint `json:"kitchen_id" validate:"required"`
}

type SubscribeKitchenAsAdminPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type OTP struct {
	Otp   string `json:"otp" validate:"required,min=6,max=6"`
	Email string `json:"email" validate:"required,email"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// go automatically maps kitchenMember to kitchen_members table
type KitchenMember struct {
	UserID    uint   `json:"user_id"`
	KitchenID uint   `json:"kitchen_id"`
	Role      string `json:"role"`
}

type KitchenMemberClaims struct {
	KitchenID uint   `json:"kitchen_id"`
	Role      string `json:"role"`
}

type Kitchen struct {
	gorm.Model
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required,min=3,max=20"`
	Slug string `json:"slug" validate:"required,min=3,max=20"`
}

type FoodItem struct {
	gorm.Model
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	KitchenID   uint    `json:"kitchen_id"`
	IsActive    bool    `json:"is_active"`
	IsVeg       bool    `json:"is_veg"`
	IsVegan     bool    `json:"is_vegan"`
}

type FoodItemImage struct {
	ID                  uint   `json:"id"`
	FoodItemID          uint   `json:"food_item_id"`
	ImageURL            string `json:"image_url"`
	OriginalDestination string `json:"original_destination"`
	DisplayOrder        int    `json:"display_order"`
}

type CreateFoodItemPayload struct {
	Name        string  `json:"name" validate:"required,min=3,max=20"`
	Description string  `json:"description" validate:"required,min=3,max=200"`
	Price       float64 `json:"price" validate:"required,numeric"`
	IsVeg       bool    `json:"is_veg"`
	IsVegan     bool    `json:"is_vegan"`
	IsActive    bool    `json:"is_active"`
}

type KitchenAndRoleClaims struct {
	KitchenID uint `json:"kitchen_id"`
	Role      string
}

type UpdateFoodItemPayload struct {
	Name        *string  `form:"name" validate:"min=3,max=20"`
	Description *string  `form:"description" validate:"min=3,max=200"`
	Price       *float64 `form:"price" validate:"numeric"`
	IsVeg       *bool    `form:"is_veg"`
	IsVegan     *bool    `form:"is_vegan"`
	IsActive    *bool    `form:"is_active"`

	ExistingImageIDs []string `form:"existing_image_ids"`
}

type ImageDataStruct struct {
	ImageURL            string `json:"imageUrl"`
	OriginalDestination string `json:"originalDestination"`
	OriginalOrder       int    `json:"originalOrder"`
}

type CreateMenuPayload struct {
	Date       time.Time  `json:"date" validate:"required,datetime"`
	Status     string     `json:"status" validate:"required,oneof=active"`
	OrderOpen  *time.Time `json:"order_open" validate:"datetime"`
	OrderClose *time.Time `json:"order_close" validate:"datetime"`
	Items      []int      `json:"items" validate:"required,min=1"`
}

const MenuStatusActive = "active"
const MenuStatusDraft = "draft"

type Menu struct {
	gorm.Model
	ID         uint       `json:"id"`
	Date       time.Time  `json:"date" validate:"required,datetime"`
	Status     string     `json:"status" validate:"required,oneof=active"`
	OrderOpen  *time.Time `json:"order_open" validate:"datetime"`
	OrderClose *time.Time `json:"order_close" validate:"datetime"`
	UpdatedBy  uint       `json:"updated_by" validate:"required,uint"`
	CreatedBy  uint       `json:"created_by" validate:"required,uint"`
	KitchenID  uint       `json:"kitchen_id" validate:"required,uint"`
}

type MenuItem struct {
	ID          uint    `json:"id"`
	MenuID      uint    `json:"menu_id"`
	FoodItemID  uint    `json:"food_item_id"`
	Price       float64 `json:"price" validate:"numeric"`
	IsAvailable bool    `json:"is_available"`
}

type PublishMenuPayload struct {
	Status string `json:"status" validate:"required,oneof=active"`
}

type UpdateMenuPayload struct {
	Items      []int      `json:"items" validate:"required,min=1"`
	Date       *time.Time `json:"date"`
	OrderOpen  *time.Time `json:"order_open"`
	OrderClose *time.Time `json:"order_close"`
}
