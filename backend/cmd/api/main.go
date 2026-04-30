package main

import (
	"fmt"

	"food-serve.com/internal/db"
	"food-serve.com/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var Env config.Config = config.LoadConfig()


func main(){
  //load env
	//connect to db 
	//setup server
	//creating a router from gin
	r := gin.Default()

	dsn := (&mysql.Config{
		User: Env.DBUser,
		Passwd: Env.DBPass,
		Addr: Env.DBAddr,
		DBName: Env.DBName,
	}).FormatDSN()

	database, err := db.NewDBStorage(dsn)

	if err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("db connected wohooo!!")

	fmt.Println(database)

	//now we initialize into services and stores

	//start server 
	r.Run(":8080")
}