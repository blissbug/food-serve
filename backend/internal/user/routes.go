package user

import "github.com/gin-gonic/gin"

func (userService *Handler) RegisterRoutes(s *gin.Engine) {
	s.POST("/login", userService.HandleLogin)
	s.POST("/register", userService.HandleRegister)
}
