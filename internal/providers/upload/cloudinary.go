package upload

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Cloudinary struct {
	APIKey           string
	APISecret        string
	CloudName        string
	URL              string
	HttpClient       *http.Client
	CloudinaryClient *cloudinary.Cloudinary
}

func NewCloudinary(apiKey, apiSecret, cloudName, url string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &Cloudinary{
		APIKey:           apiKey,
		APISecret:        apiSecret,
		CloudName:        cloudName,
		URL:              url,
		HttpClient:       &http.Client{Timeout: 15 * time.Second},
		CloudinaryClient: cld,
	}, nil
}

func (c *Cloudinary) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	resp, err := c.CloudinaryClient.Upload.Upload(ctx, file, uploader.UploadParams{
		UniqueFilename: api.Bool(true),
		Folder:         "mekoko",
		PublicID:       header.Filename,
		Overwrite:      api.Bool(false),
		ResourceType:   api.Image.String(),
		AssetFolder:    "mekoko",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
