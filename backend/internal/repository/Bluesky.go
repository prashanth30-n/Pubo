package repository

import (
	"context"
	
	"fmt"

	"github.com/PatibandlaVenkat/Pubo/internal/server"
	
	"github.com/jackc/pgx/v5"
	"github.com/PatibandlaVenkat/Pubo/internal/model/Bluesky"

)
type BlueskyRepository struct{
	server *server.Server
}
func NewBlueskyRepository(server *server.Server) *BlueskyRepository{
	return &BlueskyRepository{
		server:server,
	}
}

func(r *BlueskyRepository) ConnectAccounts(ctx context.Context,account *bluesky.ConnectedAccounts)(*bluesky.ConnectedAccounts,error){
	stmt:=`INSERT INTO connected_accounts (
    user_id,
    handle,
    platform_id,
    access_token_encrypted,
    refresh_token_encrypted,
    display_name,
    password,
    did,
    pds_url
)
VALUES (
    @user_id,
    @handle,
    @platform_id,
    @access_token_encrypted,
    @refresh_token_encrypted,
    @display_name,
    @password,
    @did,
    @pds_url
)
ON CONFLICT (user_id, platform_id)
DO UPDATE SET
    handle = EXCLUDED.handle,
    access_token_encrypted = EXCLUDED.access_token_encrypted,
    refresh_token_encrypted = EXCLUDED.refresh_token_encrypted,
    display_name = EXCLUDED.display_name,
    password = EXCLUDED.password,
    did = EXCLUDED.did,
    pds_url = EXCLUDED.pds_url
RETURNING *`;
	   rows,err:=r.server.DB.Pool.Query(ctx,stmt,pgx.NamedArgs{
		"user_id":account.UserID,
		"handle":account.Handle,
		"platform_id":account.PlatformId,
		"access_token_encrypted":account.AccessToken,
		"refresh_token_encrypted":account.RefreshToken,
		"display_name":account.DisplayName,
		"password":account.Password,
		"did":account.DID,
		"pds_url":account.PDSURL,
	   })
	   if err!=nil{
		return nil,fmt.Errorf("failed to execute created bluesky connected account",err)
	   }
	   BlueskyItem,err:=pgx.CollectOneRow(rows,pgx.RowToStructByName[bluesky.ConnectedAccounts])
	   if err!=nil{
		return nil,fmt.Errorf("failed to collect a row from table %w",err)
	
	   }
	
	   return &BlueskyItem,nil

}
