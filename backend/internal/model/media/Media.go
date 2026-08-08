package media

import (
	"time"

	"github.com/prashanth30-n/Pubo/internal/model"

	"github.com/google/uuid"
)

type Asset struct {
	model.Base
	ID            uuid.UUID `json:"id"`
	OwnerUserID   string    `json:"ownerUserId"`
	OriginalName  string    `json:"originalName"`
	BlobName      string    `json:"blobName"`
	StorageURL    string    `json:"storageURL"`
	ContainerName string    `json:"containerName"`
	ContentType   string    `json:"contentType"`
	SizeBytes     int64     `json:"sizeBytes"`
	ETag          string    `json:"etag"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
