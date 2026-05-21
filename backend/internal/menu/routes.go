package menu

import (
	"food-serve.com/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (h Handler) RegisterRoutes(r *gin.RouterGroup) {
	kitchenGroup := r.Group("/kitchens/:kitchenId")

	menuGroup := kitchenGroup.Group("/menu")
	//accessible to anyone subscribed to kitchen if published
	menuGroup.GET("/today", h.GetMenuHandler)

	adminOnlyMenuGroup := menuGroup.Use(middleware.IsAdmin())
	// add middleware to check if it is the admin else restrict
	//kitchen admin only
	adminOnlyMenuGroup.POST("/create", h.CreateMenuHandler)
	adminOnlyMenuGroup.PATCH(":menuId/publish", h.PublishMenuHandler)
	adminOnlyMenuGroup.PATCH("/:menuId/update", h.UpdateMenuHandler)
	adminOnlyMenuGroup.DELETE("/:menuId/delete", h.DeleteMenuHandler)
}
