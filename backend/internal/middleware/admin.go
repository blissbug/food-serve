package middleware

import (
	"errors"
	"strconv"

	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
)

func IsAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value := ctx.Param("kitchenId")
		kitchenId, err := strconv.Atoi(value)

		if err != nil {
			response.Error(ctx, errors.New("invalid kitchen id"), 400)
			ctx.Abort()
			return
		}

		initialValue, exist := ctx.Get(KitchenAndRoleKey)

		if !exist {
			response.Error(ctx, errors.New("no kitchen data found, please subscribe or create a kitchen"), 400)
			ctx.Abort()
		}

		KitchensAndRoles, ok := initialValue.([]types.KitchenAndRoleClaims)

		if !ok {
			response.Error(ctx, err, 400)
			ctx.Abort()
			return
		}

		isAdminOfThisKitchen := utils.CheckRoleForThisKitchen(uint(kitchenId), "admin", KitchensAndRoles)

		if !isAdminOfThisKitchen {
			response.Error(ctx, err, 400)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
