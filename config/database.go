package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*

 THIS FILE IS MAINLY USED TO CONNECT DATABASES (MYSQL,MONGODB) FOR THIS APPLICATION

*/

var Mysql *sql.DB
var MongoDBClient *mongo.Client
var MongoDB *mongo.Database

// Connecting mysql
func ConnectMYSQL() {

	var connectionErr error

	Mysql, connectionErr = sql.Open("mysql", os.Getenv("MYSQL_DB_URL"))

	if connectionErr != nil {
		fmt.Println("MYSQL connection error : ", connectionErr)
		return
	}

	fmt.Println("Mysql connected !!")

}

// Connecting mongodb
func ConnectMongoDB() {

	//Creating a context with 10 seconds
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
	client, clientErr := mongo.Connect(context, clientOptions)

	if clientErr != nil {
		fmt.Println("Client error : ", clientErr)
		return
	}

	pingErr := client.Ping(context, nil)

	if pingErr != nil {
		fmt.Println("Ping error : ", pingErr)
		return
	}

	MongoDBClient = client
	MongoDB = client.Database("chitchat")

	fmt.Println("Mongodb connected !!")
}
