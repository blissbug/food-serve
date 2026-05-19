package asset_upload

import "github.com/cloudinary/cloudinary-go/v2"

func InitConfig(cldUrl string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(cldUrl)

	if err != nil {
		return cld, err
	}

	return cld, nil
}
