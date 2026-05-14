package user

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterRoutes(s *gin.Engine) {
	s.POST("/login", h.HandleLogin)
	s.POST("/register", h.HandleRegister)
}
