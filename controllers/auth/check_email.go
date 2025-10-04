package auth

import (
	"chitchat/config"
	"chitchat/helpers"
	authmodel "chitchat/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ------- CHECK USERNAME AND EMAIL CONTROLLER ------------------
// To check if the username or email already presented
func CheckEmailController(ctx *gin.Context) {
	var checkModel authmodel.CheckUsernameAndEmailModel

	//Parsing the data from json
	parseErr := ctx.ShouldBind(&checkModel)

	if parseErr != nil {
		fmt.Println("Parsing error : ", parseErr)
		helpers.SendMessageAsJson(ctx, "Something went wrong", http.StatusInternalServerError)
		return
	}

	//Counting email from table
	var emailCount int64

	emailCountErr := config.GORM.
		Table("userdetails").
		Where("email=?", checkModel.Email).
		Count(&emailCount).Error

	if emailCountErr != nil {
		fmt.Println("Counting error /check", emailCountErr)
		return
	}

	//Checking whether the emailcount is 1 or 0
	if int(emailCount) != 0 {
		helpers.SendMessageAsJson(ctx, "Email is already in use", http.StatusAlreadyReported)
		return
	}

	helpers.SendMessageAsJson(ctx, "Success", http.StatusOK)
}
