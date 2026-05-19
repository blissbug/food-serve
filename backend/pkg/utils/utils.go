package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"food-serve.com/pkg/config"
	"food-serve.com/pkg/types"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPasswordInBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedPasswordInBytes)

	return hashedPassword, nil
}

func ComparePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func JWTwithClaimsForTokens(user types.User, kitchenMember []types.KitchenMember) (string, string, error) {
	Env, err := config.LoadConfig()
	fmt.Println(Env)

	if err != nil {
		return "", "", err
	}

	var accessToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":        user.ID,
		"email":         user.Email,
		"kitchenMember": kitchenMember,
		"exp":           time.Now().Add(time.Minute * 60).Unix(),
		"iat":           time.Now().Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte(Env.JWTSecret))

	var refreshToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), //goes for a week
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(Env.JWTSecret))

	return accessTokenString, refreshTokenString, err
}

func CheckRoleForThisKitchen(kitchenId uint, role string, kitchenMembers []types.KitchenAndRoleClaims) bool {
	for i := 0; i < len(kitchenMembers); i++ {
		kitchenMember := kitchenMembers[i]
		if kitchenMember.KitchenID == kitchenId && kitchenMember.Role == role {
			return true
		}
	}
	return false
}

func UploadImagesToCloudinary(files []*multipart.FileHeader, cld *cloudinary.Cloudinary, ctx *gin.Context) ([]types.ImageDataStruct, error) {
	var ItemImageData []types.ImageDataStruct
	for order, file := range files {
		err := os.MkdirAll("./assets/uploads", os.ModePerm)
		if err != nil {
			return ItemImageData, err
		}
		// Upload the file to specific dst.
		dst := filepath.Join("./assets/uploads", fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(file.Filename)))
		fileSaveError := ctx.SaveUploadedFile(file, dst)

		if fileSaveError != nil {
			return ItemImageData, fileSaveError
		}

		//cloudinary job
		resp, uploadErr := cld.Upload.Upload(ctx, dst, uploader.UploadParams{PublicID: uuid.NewString()})

		if uploadErr != nil {
			return ItemImageData, uploadErr
		}
		//save with old file name
		fmt.Println(resp.URL)

		ItemImageData = append(ItemImageData, types.ImageDataStruct{
			OriginalDestination: dst,
			ImageURL:            resp.URL,
			OriginalOrder:       order,
		})

		//delete it from local storage
		fileRemoveErr := os.Remove(dst)
		if fileRemoveErr != nil {
			return ItemImageData, fileRemoveErr
		}
	}
	return ItemImageData, nil
}
