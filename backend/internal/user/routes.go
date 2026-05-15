package user

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(s *gin.Engine) {
	s.POST("/login/send-otp", h.HandleLogin)
	s.POST("/login/verify", h.HandlerVerifyLogin)
	s.POST("/auth/refresh", h.HandleRefresh)
	s.POST("/register/send-otp", h.HandleRegister)
	s.POST("/register/verify", h.HandleVerifyRegistration)
}
