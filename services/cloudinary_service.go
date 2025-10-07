package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

/*

	THIS FILE IS MAINLY USED FOR UPLOADING FILE WITH URL AND FILE IN FORMDATA

*/

// ---------- UPLOADING FILE WITH URL -----------------
func UploadFileWithURL(fileUrl string) string {
	//Forming cloudinary
	cld, cldErr := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if cldErr != nil {
		fmt.Println("Cloudinary forming error : ", cldErr)
		return ""
	}
	//Creating a context
	context := context.Background()

	if fileUrl != "" {
		response, responseErr := cld.Upload.Upload(context, fileUrl, uploader.UploadParams{PublicID: uuid.NewString()})

		if responseErr != nil {
			fmt.Println("Response error : ", responseErr)
			return ""
		}

		return response.SecureURL
	}
	return ""

}

// ------------ FILE UPLOADING IN FORMDATA-------------------
func UploadFileInFormData(file *multipart.FileHeader) (string, string) {

	if file == nil {
		return "", ""
	}
	//Forming cloudinary
	cld, cldErr := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if cldErr != nil {
		fmt.Println("Cloudinary forming error : ", cldErr)
		return "", ""
	}
	//Creating context
	context := context.Background()
	//Opening the file
	source, openErr := file.Open()

	if openErr != nil {
		fmt.Println("Opening error : ", openErr)
		return "", ""
	}

	publicId := uuid.NewString()

	response, responseErr := cld.Upload.Upload(context, source, uploader.UploadParams{PublicID: publicId})

	if responseErr != nil {
		fmt.Println("Response error : ", responseErr)
		return "", ""
	}

	return response.SecureURL, publicId
}

// ----------- DELETE FILE -----------
// To delete any file using the public id of the file
func DeleteFile(publicId string) error {
	cld, cldErr := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if cldErr != nil {
		return cldErr
	}

	response, responseErr := cld.Upload.Destroy(context.TODO(), uploader.DestroyParams{PublicID: publicId})

	if responseErr != nil {
		return responseErr
	}

	fmt.Println("Response from cloudinary : ", response.Response)

	return nil
}
