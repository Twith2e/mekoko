package upload

import (
	"context"
	"log"
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
	HttpClient       *http.Client
	CloudinaryClient *cloudinary.Cloudinary
}

func NewCloudinary(apiKey, apiSecret, cloudName string) (*Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("cloudinary new from params: failed to create cloudinary client: %s\n", err)
		return nil, err
	}
	return &Cloudinary{
		APIKey:           apiKey,
		APISecret:        apiSecret,
		CloudName:        cloudName,
		HttpClient:       &http.Client{Timeout: 15 * time.Second},
		CloudinaryClient: cld,
	}, nil
}

func (c *Cloudinary) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// api.SetTraceLevel(api.TraceLevel("debug"))
	resp, err := c.CloudinaryClient.Upload.Upload(ctx, file, uploader.UploadParams{
		UniqueFilename: api.Bool(true),
		Folder:         "mekoko",
		PublicID:       header.Filename,
		Overwrite:      api.Bool(false),
		ResourceType:   api.Image.String(),
	})
	if err != nil {
		log.Printf("cloudinary upload file: failed to upload file: %s\n", err)
		return "", err
	}
	log.Printf("cloudinary response: public_id=%s, url=%s, secure_url=%s\n", resp.PublicID, resp.URL, resp.SecureURL)
	return resp.SecureURL, nil
}
