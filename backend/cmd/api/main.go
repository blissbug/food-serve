package main

import (
	"fmt"

	cache "food-serve.com/internal/cache"
	"food-serve.com/internal/db"
	"food-serve.com/internal/user"
	"food-serve.com/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

	r := gin.Default()

	dsn := (&mysql.Config{
		User:                 Env.DBUser,
		Passwd:               Env.DBPass,
		Net:                  "tcp",
		Addr:                 Env.DBAddr,
		DBName:               Env.DBName,
		AllowNativePasswords: true,
	}).FormatDSN()

	database, err := db.NewDBStorage(dsn)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("db connected wohooo!!")

	//SETUP REDIS SERVER - run container before this pls
	cache.RedisClient(Env)

	userStore := user.NewStore(database) //now userStore has db inside it
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(r)

	//now we initialize into services and stores
	//start server
	err = r.Run(":8080")

	if err != nil {
		fmt.Println(err)
	}
}
