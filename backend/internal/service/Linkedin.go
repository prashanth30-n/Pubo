package service

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/prashanth30-n/Pubo/internal/middleware"
	linkedin "github.com/prashanth30-n/Pubo/internal/model/Linkedin"
	"github.com/prashanth30-n/Pubo/internal/repository"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type LinkedinService struct {
	server            *server.Server
	LikedinRepository *repository.LinkedinRepository
	LinkedinClient    *LinkedinClient
}

func NewLinkedinService(s *server.Server, r *repository.LinkedinRepository) *LinkedinService {
	return &LinkedinService{
		server:            s,
		LikedinRepository: r,
		LinkedinClient:    NewLinkedinClient(),
	}

}
func (s *LinkedinService) CreateLinkedinConnection(ctx echo.Context, userId string, payload *linkedin.CreateLinkedinConnectionPayload) (*linkedin.LinkedinConnectedAccount, error) {
	logger := middleware.GetLogger(ctx)
	UserInfo, err := s.LinkedinClient.UserInfo(ctx.Request().Context(), payload.AccessToken)
	if err != nil {
		logger.Error().Err(err).Msg("Linkedin token validation failed")
		return nil, fmt.Errorf("invalid linkedin access token: %w", err)
	}
	var avatarURL *string
	if UserInfo.Picture != "" {
		avatarURL = &UserInfo.Picture
	}
	displayName := payload.DisplayName
	account := &linkedin.LinkedinConnectedAccount{
		UserID:      userId,
		PlatformId:  2,
		DisplayName: displayName,
		AvatarUrl:   avatarURL,
		AccessToken: payload.AccessToken,
		DID:         UserInfo.Sub,
		PDSURL:      "https://api.linkedin.com",
		Is_active:   true,
		// TokenExpiry:time.Now().Add(60*24*time.Hour),
	}

	LinkedinItem, err := s.LikedinRepository.ConnectLinkedinAccount(ctx.Request().Context(), account)
	if err != nil {
		logger.Error().Err(err).Msg("failed to connect Linkedin account")
		return nil, err
	}
	logger.Info().Str("event", "linkedin_account_connected").Str("member_id", UserInfo.Sub).Msg("Linkedin account sucessfully connected")
	return LinkedinItem, nil
}
