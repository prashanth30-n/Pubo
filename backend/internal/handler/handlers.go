package handler

import (
	"github.com/prashanth30-n/Pubo/internal/server"
	"github.com/prashanth30-n/Pubo/internal/service"
)

type Handlers struct {
	Health   *HealthHandler
	OpenAPI  *OpenAPIHandler
	Quotes   *QuoteHandler
	Media    *MediaHandler
	Bluesky  *BlueskyHandler
	Signup   *SignUpHandler
	Linkedin *LinkedinHandler
	Posts    *PostHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:   NewHealthHandler(s),
		OpenAPI:  NewOpenAPIHanlder(s),
		Quotes:   NewQuoteHandler(s, services.Quotes),
		Media:    NewMediaHandler(s, services.MediaService),
		Bluesky:  NewBlueskyHandler(s, services.BlueskyService),
		Signup:   NewSignUpHandler(s, services.SignUpService),
		Linkedin: NewLinkedinHandler(s, services.LinkedinService),
		Posts:    NewPostHandler(s, services.PostService),
	}
}
