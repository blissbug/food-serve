package types

type User struct {
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
	ID   uint   `json:"id"`
	Name string `json:"name" validate:"required,min=3,max=20"`
	Slug string `json:"slug" validate:"required,min=3,max=20"`
}

type FoodItem struct {
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
	FoodItemID          uint   `json:"food_item_id"`
	ImageURL            string `json:"image_url"`
	OriginalDestination string `json:"original_destination"`
	DisplayOrder        int    `json:"display_order"`
}

type CreateFoodItemPayload struct {
	Name        string  `json:"name" validate:"required,min=3,max=20"`
	Description string  `json:"description" validate:"required,min=3,max=200"`
	Price       float64 `json:"price" validate:"required,numeric"`
	KitchenID   uint    `json:"kitchen_id" validate:"required"`
	IsVeg       bool    `json:"is_veg"`
	IsVegan     bool    `json:"is_vegan"`
	IsActive    bool    `json:"is_active"`
}

type KitchenAndRoleClaims struct {
	KitchenID uint `json:"kitchen_id"`
	Role      string
}
