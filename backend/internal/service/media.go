package service

import (
	"errors"
	"mime/multipart"
	"strings"

	"github.com/prashanth30-n/Pubo/internal/lib/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	"github.com/prashanth30-n/Pubo/internal/model/media"
	"github.com/prashanth30-n/Pubo/internal/repository"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type MediaService struct {
	server    *server.Server
	MediaRepo repository.MediaRepository
	blob      *storage.AzureBlobClient
}

func NewMediaService(s *server.Server, MediaRepo repository.MediaRepository, blob *storage.AzureBlobClient) *MediaService {
	return &MediaService{
		server:    s,
		MediaRepo: MediaRepo,
		blob:      blob,
	}
}
func (s *MediaService) Upload(ctx echo.Context, ownerUserID string, fileHeader *multipart.FileHeader, folder string) (*media.Asset, error) {
	if fileHeader == nil {
		return nil, errors.New("file is required")
	}
	logger := middleware.GetLogger(ctx)
	ct := strings.ToLower(fileHeader.Header.Get("Content-Type"))
	if !strings.HasPrefix(ct, "image/") && !strings.HasPrefix(ct, "video/") && !strings.HasPrefix(ct, "application/") {
		return nil, errors.New("unsupported file type")
	}
	uploadResult, err := s.blob.UploadFile(ctx.Request().Context(), fileHeader, folder)
	if err != nil {
		logger.Error().Err(err).Msg("failed to upload the image")
		return nil, err
	}
	id, err := uuid.NewUUID()
	asset := media.Asset{
		ID:            id,
		OwnerUserID:   ownerUserID,
		OriginalName:  fileHeader.Filename,
		BlobName:      uploadResult.BlobName,
		StorageURL:    uploadResult.StorageURL,
		ContainerName: uploadResult.ContainerName,
		ContentType:   uploadResult.ContentType,
		SizeBytes:     uploadResult.SizeBytes,
		ETag:          uploadResult.ETag,
	}
	if err := s.MediaRepo.Create(ctx.Request().Context(), asset); err != nil {
		_ = s.blob.DeleteBlob(ctx.Request().Context(), uploadResult.BlobName)

		logger.Error().Err(err).Msg("failed to upload the image to DB")
		return nil, err

	}
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().Str("event", "image_uploaded").
		Msg("file uploaded sucessfully")

	return &asset, nil
}
func (s *MediaService) Get(ctx echo.Context, id string) (*media.Asset, error) {
	var asset *media.Asset
	logger := middleware.GetLogger(ctx)
	asset, err := s.MediaRepo.GetByID(ctx.Request().Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to retrieve images by id")
		return nil, err
	}
	return asset, nil

}
func (s *MediaService) List(ctx echo.Context, ownerUserID string, limit int, offset int) ([]media.Asset, error) {
	var asset []media.Asset
	logger := middleware.GetLogger(ctx)
	asset, err := s.MediaRepo.ListByOwner(ctx.Request().Context(), ownerUserID, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch images by ownerUserID")
		return nil, err
	}
	return asset, nil

}
func (s *MediaService) Delete(ctx echo.Context, id string) error {
	logger := middleware.GetLogger(ctx)
	asset, err := s.MediaRepo.GetByID(ctx.Request().Context(), id)
	if err != nil {

		return err
	}
	if asset == nil {
		return errors.New("media not found")
	}
	if err := s.blob.DeleteBlob(ctx.Request().Context(), asset.BlobName); err != nil {
		logger.Error().Err(err).Msg("failed to delete the file in blob")
		return err
	}
	err = s.MediaRepo.Delete(ctx.Request().Context(), id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delte the file from db")
		return err
	}
	return nil

}
