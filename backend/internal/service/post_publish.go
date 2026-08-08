package service

import (
	"context"
	"fmt"

	"github.com/prashanth30-n/Pubo/internal/model/media"
	"github.com/prashanth30-n/Pubo/internal/repository"
)

func (s *PostService) publishBluesky(ctx context.Context, a *repository.ConnectedAccount, text string, assets []media.Asset) (string, error) {
	if len([]rune(text)) > 300 {
		return "", fmt.Errorf("bluesky posts are limited to 300 characters")
	}
	client := NewBskyClient()
	if a.PDSURL != "" {
		client.PDSURL = a.PDSURL
	}
	images := []BlueskyImage{}
	for _, asset := range assets {
		data, err := s.storage.Download(ctx, asset.BlobName)
		if err != nil {
			return "", err
		}
		blob, err := client.UploadBlob(ctx, a.AccessToken, data, asset.ContentType)
		if err != nil {
			return "", err
		}
		images = append(images, BlueskyImage{Alt: "", Image: blob})
	}
	return client.CreatePost(ctx, a.AccessToken, a.DID, text, images)
}
func (s *PostService) publishLinkedin(ctx context.Context, a *repository.ConnectedAccount, text string, assets []media.Asset) (string, error) {
	if len([]rune(text)) > 3000 {
		return "", fmt.Errorf("linkedin posts are limited to 3000 characters")
	}
	data := [][]byte{}
	types := []string{}
	for _, asset := range assets {
		image, err := s.storage.Download(ctx, asset.BlobName)
		if err != nil {
			return "", err
		}
		data = append(data, image)
		types = append(types, asset.ContentType)
	}
	return s.linkedin.CreatePost(ctx, a.AccessToken, "urn:li:person:"+a.DID, text, data, types)
}
