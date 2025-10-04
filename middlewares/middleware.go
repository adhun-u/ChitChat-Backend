package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*

 MAINLY USING FOR CHECKING WHETHER THE REQUEST HAS TOKEN OR NOT

*/

func MiddleWare(ctx *gin.Context) {

	//Parsing the token from headers
	headerToken := ctx.GetHeader("Authorization")

	if headerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Token is required"})
		return
	}

	//Trimming white space infront of the token and endof the token
	token := strings.Trim(headerToken, " ")

	//Parsing the token
	parsedToken, parseErr := jwt.Parse(token, func(jwtToken *jwt.Token) (any, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method : %s", jwtToken.Header)
		}
		return []byte(os.Getenv("JWT_SECRETE")), nil
	})

	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid token"})
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid token"})
		return
	}

	//Checking whether the token contains authtype and userid
	currentUserid, containsId := claims["userId"]
	_, containsAuthtype := claims["authtype"]

	if containsId && containsAuthtype {
		//Setting the current userid
		ctx.Set("userId", currentUserid)
		//Moving to next context
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid token"})
	}

}
