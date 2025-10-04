package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
	THIS FILE IS MAINLY USED TO INITIALIZE GORM WITH MYSQL DATABASE
*/

var GORM *gorm.DB

// Initializing gorm with mysql
func InitGORM() {
	var err error
	fmt.Println()
	GORM, err = gorm.Open(mysql.New(mysql.Config{
		Conn: Mysql,
	}), &gorm.Config{})
	if err != nil {
		fmt.Println("Could not initialize gorm", err)
		return
	}

}
