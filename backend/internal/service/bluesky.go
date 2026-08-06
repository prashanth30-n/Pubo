package service

import (
	"fmt"

	"github.com/PatibandlaVenkat/Pubo/internal/middleware"
	bluesky "github.com/PatibandlaVenkat/Pubo/internal/model/Bluesky"
	"github.com/PatibandlaVenkat/Pubo/internal/repository"
	"github.com/PatibandlaVenkat/Pubo/internal/server"
	"github.com/labstack/echo/v4"
)
type BlueskyService struct{
	server *server.Server
	blueskyrepo *repository.BlueskyRepository
	bskyClient *BskyClient
}
func NewBlueskyService(s *server.Server,blueskyrepo *repository.BlueskyRepository) (*BlueskyService){
	return &BlueskyService{
		server:s,
		blueskyrepo: blueskyrepo,
		bskyClient: NewBskyClient(),

	}
}
func (s*BlueskyService) ConnectBluesky(ctx echo.Context,userID string,payload *bluesky.CreateBlueskyConnectionPayload)(*bluesky.ConnectedAccounts,error){
	logger:=middleware.GetLogger(ctx)
	session,err:=s.bskyClient.CreateSession(ctx.Request().Context(),*payload.Handle,*payload.Password)
	if err!=nil{
		logger.Error().Err(err).Str("handle",*payload.Handle).Msg("bluesky createSession failed")
		return nil,fmt.Errorf("invalid bluesky credentials or server error: %w",err)
	}
	account:=&bluesky.ConnectedAccounts{
		UserID: userID,
		Handle: session.Handle,
		DID:session.DID,
		AccessToken:session.AccessJwt,
		RefreshToken:session.RefreshJwt,
		PDSURL: s.bskyClient.PDSURL,
		PlatformId: payload.PlatformId,
		DisplayName: payload.DisplayName,
		
	}
	BlueskyItem,err:=s.blueskyrepo.ConnectAccounts(ctx.Request().Context(),account)
	if err!=nil{
		logger.Error().Err(err).Msg("failed to connect bluesky account")
		return nil,err
	}
	//bussiness event logging
	eventLogger:=middleware.GetLogger(ctx)
	eventLogger.Info().Str("event","bluesky_account_connected").Msg("Bluesky account successfully created")
	
return BlueskyItem,nil
}
