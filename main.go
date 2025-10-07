package main

import (
	"chitchat/config"
	"chitchat/services"
	"chitchat/urls"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadErr := godotenv.Load()

	if loadErr != nil {
		fmt.Println("ENV loading error : ", loadErr)
		return
	}
	//Connecting MYSQL database
	config.ConnectMYSQL()
	//Closing mysql connection when shutting down
	defer config.Mysql.Close()
	//Connecting mongodb
	config.ConnectMongoDB()
	//Initializing gorm
	config.InitGORM()
	//Initializing route
	route := gin.Default()
	//Registering necessary urls
	urls.RegisterAuthUrls(route)
	urls.RegisterUserUrls(route)
	urls.RegisterMessageUrls(route)
	urls.RegisterGroupUrls(route)
	//Initializing firebase
	services.InitializeFirebaseApp()
	route.Run("0.0.0.0:8080")
}
