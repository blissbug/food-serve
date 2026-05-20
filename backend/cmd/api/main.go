package main

import (
	"fmt"

	"food-serve.com/internal/app"
	cache "food-serve.com/internal/cache"
	"food-serve.com/internal/db"
	asset_upload "food-serve.com/pkg/asset-upload"
	"food-serve.com/pkg/config"
	logs "food-serve.com/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func main() {
	//load env
	//connect to db
	//setup server
	//creating a router from gin
	Env, err := config.LoadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	//validate env vars
	err = Env.Validate()

	if err != nil {
		fmt.Println(err)
		return
	}

	//setup zap logger and make it globally available
	logger := logs.InitLogger()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
	zap.L().Info("Logger was initialized!")

	r := gin.Default()

	dsn := (&mysql.Config{
		User:                 Env.DBUser,
		Passwd:               Env.DBPass,
		Net:                  "tcp",
		Addr:                 Env.DBAddr,
		DBName:               Env.DBName,
		ParseTime:            true,
		AllowNativePasswords: true,
	}).FormatDSN()

	database, err := db.NewDBStorage(dsn)

	if err != nil {
		zap.L().Error("DB connection failed", zap.Error(err))
		fmt.Println(err)
		return
	}

	fmt.Println("db connected wohooo!!")

	cld, cldErr := asset_upload.InitConfig(Env.CloudinaryUrl)

	if cldErr != nil {
		zap.L().Error("cloudinary connection failed", zap.Error(cldErr))
		fmt.Println(cldErr)
		return
	}

	//SETUP REDIS SERVER - run container before this pls
	rdb := cache.RedisClient(Env)

	//handles routes and services
	app.NewApplication(database, r, rdb, cld)

	//now we initialize into services and stores
	//start server
	err = r.Run(":8080")

	if err != nil {
		fmt.Println(err)
	}
}
