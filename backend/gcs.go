package backend

import (
	"context"
	"fmt"
	"io"

	"around/util"

	"cloud.google.com/go/storage"
)

/* How to use GCS library: https://github.com/GoogleCloudPlatform/golang-samples/blob/main/storage/objects/upload_file.go */

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

func InitGCSBackend(config *util.GCSInfo) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}

	GCSBackend = &GoogleCloudStorageBackend{
		client: client,
		bucket: config.Bucket,
	}
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	// r = content to be uploaded
	// string = name of content that is uploaded
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)

	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Everyone can read the content, but not edit it
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)

	return attrs.MediaLink, nil
}
