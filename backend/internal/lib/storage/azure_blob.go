package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/prashanth30-n/Pubo/internal/config"
)

type UploadResult struct {
	BlobName      string
	StorageURL    string
	ContainerName string
	ContentType   string
	SizeBytes     int64
	ETag          string
}
type AzureBlobClient struct {
	client        *azblob.Client
	containerName string
	accountName   string
}

func NewAzureBlobClient(cfg *config.Config) (*AzureBlobClient, error) {
	client, err := azblob.NewClientFromConnectionString(cfg.Azure.ConnectionString, nil)
	if err != nil {
		return nil, err
	}
	return &AzureBlobClient{
		client:        client,
		containerName: cfg.Azure.ContainerName,
		accountName:   cfg.Azure.AccountName,
	}, nil
}
func (a *AzureBlobClient) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*UploadResult, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octed-stream"
	}
	blobName := BuildBlobName(folder, fileHeader.Filename)
	_, err = a.client.UploadStream(ctx, a.containerName, blobName, bytes.NewReader(data), &azblob.UploadStreamOptions{
		// HTTPHeaders: &azblob.BlobHTTPHeaders{
		// 	BlobContentType:&contentType,
		// },

	})
	if err != nil {
		return nil, err
	}
	storageURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", a.accountName, a.containerName, blobName)
	return &UploadResult{
		BlobName:      blobName,
		StorageURL:    storageURL,
		ContainerName: a.containerName,
		ContentType:   contentType,
		SizeBytes:     int64(len(data)),
	}, nil

}
func (a *AzureBlobClient) DeleteBlob(ctx context.Context, blobName string) error {
	_, err := a.client.DeleteBlob(ctx, a.containerName, blobName, nil)
	return err
}
func (a *AzureBlobClient) Download(ctx context.Context, blobName string) ([]byte, error) {
	response, err := a.client.DownloadStream(ctx, a.containerName, blobName, nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
func BuildBlobName(folder, fileName string) string {
	folder = strings.Trim(folder, "/")
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	base = strings.ReplaceAll(base, " ", "_")
	ts := time.Now().UTC().Format("20060102150405")

	if folder == "" {
		return fmt.Sprintf("%s_%s%s", ts, base, ext)
	}
	return fmt.Sprintf("%s/%s_%s%s", folder, ts, base, ext)
}
