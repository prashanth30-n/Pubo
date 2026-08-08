package repository

import "github.com/prashanth30-n/Pubo/internal/server"

type Repositories struct {
	Media    *MediaRepository
	Bluesky  *BlueskyRepository
	Signup   *SignupRepository
	Linkedin *LinkedinRepository
	Posts    *PostRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Media:    NewMediaRepository(s),
		Bluesky:  NewBlueskyRepository(s),
		Signup:   NewSignupRepository(s),
		Linkedin: NewLinkedinRepository(s),
		Posts:    NewPostRepository(s),
	}
}
