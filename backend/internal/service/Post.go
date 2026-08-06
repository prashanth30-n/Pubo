package service

import (
	"github.com/PatibandlaVenkat/Pubo/internal/repository"
	"github.com/PatibandlaVenkat/Pubo/internal/server"
)
type ImageFile struct{
	Name string
	Data []byte
	ContentType string
	Size        int64

}
type PostService struct{
	server *server.Server
	PostRepo *repository.PostRepository
	BlueskyaccountRepo *repository.BlueskyRepository
	LinkedinAccountRepo *repository.LinkedinRepository
	LinkedinClient *LinkedinClient
	blueskyClient *BskyClient
}
func NewPostService(s*server.Server,pr *repository.PostRepository,)