package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/prashanth30-n/Pubo/internal/model/media"
	"github.com/prashanth30-n/Pubo/internal/repository"
)

// PublishRequest is intentionally media-ID based: the client uploads once and
// the server, not a public URL, reads the owner's assets when publishing.
type PublishRequest struct {
	Content   string   `json:"content"`
	Platforms []string `json:"platforms"`
	MediaIDs  []string `json:"mediaIds"`
}

type PlatformResult struct {
	Platform       string `json:"platform"`
	Status         string `json:"status"`
	ExternalPostID string `json:"externalPostId,omitempty"`
	Error          string `json:"error,omitempty"`
}
type PublishResponse struct {
	Results []PlatformResult `json:"results"`
}
type DraftRequest struct {
	Content  string   `json:"content"`
	MediaIDs []string `json:"mediaIds"`
}

type PostService struct {
	posts    *repository.PostRepository
	media    *repository.MediaRepository
	linkedin *LinkedinClient
	bluesky  *BskyClient
	storage  MediaReader
}

type MediaReader interface {
	Download(context.Context, string) ([]byte, error)
}

func NewPostService(posts *repository.PostRepository, mediaRepo *repository.MediaRepository, storage MediaReader) *PostService {
	return &PostService{posts: posts, media: mediaRepo, storage: storage, linkedin: NewLinkedinClient(), bluesky: NewBskyClient()}
}

func (s *PostService) Publish(ctx context.Context, userID string, req PublishRequest) (*PublishResponse, error) {
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" || len(req.Platforms) == 0 {
		return nil, fmt.Errorf("content and at least one platform are required")
	}
	assets := make([]media.Asset, 0, len(req.MediaIDs))
	for _, id := range req.MediaIDs {
		asset, err := s.media.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if asset == nil || asset.OwnerUserID != userID {
			return nil, fmt.Errorf("media %s was not found", id)
		}
		if !strings.HasPrefix(asset.ContentType, "image/") {
			return nil, fmt.Errorf("media %s is not an image", id)
		}
		assets = append(assets, *asset)
	}
	response := &PublishResponse{}
	for _, platform := range req.Platforms {
		account, err := s.posts.GetActiveAccount(ctx, userID, platform)
		if err != nil {
			response.Results = append(response.Results, PlatformResult{Platform: platform, Status: "failed", Error: err.Error()})
			continue
		}
		var externalID string
		switch platform {
		case "bluesky":
			externalID, err = s.publishBluesky(ctx, account, req.Content, assets)
		case "linkedin":
			externalID, err = s.publishLinkedin(ctx, account, req.Content, assets)
		default:
			err = fmt.Errorf("unsupported platform %q", platform)
		}
		result := PlatformResult{Platform: platform, Status: "posted", ExternalPostID: externalID}
		if err != nil {
			result.Status, result.Error = "failed", err.Error()
		}
		response.Results = append(response.Results, result)
		_ = s.posts.RecordPlatformStatus(ctx, userID, platform, req.Content, result.Status, result.ExternalPostID, result.Error)
	}
	return response, nil
}

func (s *PostService) SaveDraft(ctx context.Context, userID string, req DraftRequest) (*repository.SavedPost, error) {
	if strings.TrimSpace(req.Content) == "" && len(req.MediaIDs) == 0 {
		return nil, fmt.Errorf("a draft needs content or an image")
	}
	if _, err := s.loadImages(ctx, userID, req.MediaIDs); err != nil {
		return nil, err
	}
	return s.posts.CreateDraft(ctx, userID, req.Content, req.MediaIDs)
}
func (s *PostService) List(ctx context.Context, userID, status string, limit, offset int) (*repository.PostPage, error) {
	if status != "" && status != "draft" && status != "posted" && status != "failed" {
		return nil, fmt.Errorf("unsupported post status")
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.posts.ListPosts(ctx, userID, status, limit, offset)
}
func (s *PostService) loadImages(ctx context.Context, userID string, mediaIDs []string) ([]media.Asset, error) {
	assets := make([]media.Asset, 0, len(mediaIDs))
	for _, id := range mediaIDs {
		asset, err := s.media.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if asset == nil || asset.OwnerUserID != userID {
			return nil, fmt.Errorf("media %s was not found", id)
		}
		if !strings.HasPrefix(asset.ContentType, "image/") {
			return nil, fmt.Errorf("media %s is not an image", id)
		}
		assets = append(assets, *asset)
	}
	return assets, nil
}
