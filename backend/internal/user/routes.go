package user

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(s *gin.Engine) {
	authGroup := s.Group("/auth")
	authGroup.POST("/login/send-otp", h.HandleLogin)
	authGroup.POST("/login/verify", h.HandlerVerifyLogin)
	authGroup.POST("/refresh", h.HandleRefresh)
	authGroup.POST("/register/send-otp", h.HandleRegister)
	authGroup.POST("/register/verify", h.HandleVerifyRegistration)
}
