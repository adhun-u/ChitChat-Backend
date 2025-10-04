package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

/*

 THIS FILE IS MAINLY USED FOR GENERATING JWT TOKEN

*/

// For generating jwt token for authentication
func GenerateJWTtoken(userId int, authtype string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   userId,
		"authtype": authtype,
	})

	stringToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRETE")))

	if err != nil {
		fmt.Println("Token generation error : ", err)
		return ""
	}
	return stringToken

}
