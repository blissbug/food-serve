package kitchen

import "github.com/gin-gonic/gin"

func (h KitchenHandler) RegisterRoutes(r *gin.RouterGroup) {
	kitchenGroup := r.Group("/kitchen")
	kitchenGroup.POST("/create", h.CreateKitchen)
	kitchenGroup.POST("/enter", h.SubscribeToKitchenAsUser)
	kitchenGroup.POST("/admin/:kitchenId/enter", h.SubscribeToKitchenAsAdmin)
}
