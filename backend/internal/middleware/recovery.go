package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"food-serve.com/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var e error
				if recoveredErr, ok := err.(error); ok {
					e = recoveredErr
				} else {
					e = errors.New(fmt.Sprint(err))
				}
				stack := debug.Stack()
				zap.L().Error("Something broke!", zap.Any("error", e), zap.ByteString("stack", stack))
				response.Error(c, e, 500)
				c.Abort()
			}
		}()
		c.Next()
	}
}
