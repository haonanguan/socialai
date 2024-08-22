package backend

import (
	"context"
	"fmt"
	"io"

	"socialai/util"

	"cloud.google.com/go/storage"
)

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
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName) //place holder
	wc := object.NewWriter(ctx)                                        //copy file
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

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

func (backend *GoogleCloudStorageBackend) DeleteFromGCS(objectName string) error {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)

	// Delete the object
	err := object.Delete(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return fmt.Errorf("object %s not found in GCS", objectName)
		}
		return fmt.Errorf("failed to delete object %s from GCS: %v", objectName, err)
	}

	fmt.Printf("File %s deleted from GCS\n", objectName)
	return nil
}
