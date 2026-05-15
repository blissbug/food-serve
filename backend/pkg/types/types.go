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

type OTP struct {
	Otp   string `json:"otp" validate:"required,min=6,max=6"`
	Email string `json:"email" validate:"required,email"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
