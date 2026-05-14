package logs

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	return logger
}
