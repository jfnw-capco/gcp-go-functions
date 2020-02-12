package nozzle

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const publicURL = "https://storage.googleapis.com/%s/%s"

// UploadToBucket creates an object in a given GCP bucket with the contents of io.Reader and returns its public url
func UploadToBucket(r io.Reader, contentType string, bucketName string, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	w := bucket.Object(objectName).NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, r); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return fmt.Sprintf(publicURL, bucketName, objectName), nil
}
