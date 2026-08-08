package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	linkedin "github.com/prashanth30-n/Pubo/internal/model/Linkedin"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type LinkedinRepository struct {
	server *server.Server
}

func NewLinkedinRepository(s *server.Server) *LinkedinRepository {
	return &LinkedinRepository{
		server: s,
	}

}
func (r *LinkedinRepository) ConnectLinkedinAccount(ctx context.Context, account *linkedin.LinkedinConnectedAccount) (*linkedin.LinkedinConnectedAccount, error) {
	stmt := `INSERT INTO connected_accounts (
    user_id,
    handle,
    platform_id,
    access_token_encrypted,
    display_name,
    did,
    pds_url,
	avatar_url
)
VALUES (
    @user_id,
    @handle,
    @platform_id,
    @access_token_encrypted,
    @display_name,
    @did,
    @pds_url,
	@avatar_url
)
ON CONFLICT (user_id, platform_id)
DO UPDATE SET
    handle = EXCLUDED.handle,
    access_token_encrypted = EXCLUDED.access_token_encrypted,
    display_name = EXCLUDED.display_name,
    did = EXCLUDED.did,
    pds_url = EXCLUDED.pds_url
RETURNING *`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":                account.UserID,
		"handle":                 account.Handle,
		"platform_id":            account.PlatformId,
		"access_token_encrypted": account.AccessToken,
		"display_name":           account.DisplayName,
		"did":                    account.DID,
		"pds_url":                account.PDSURL,
		"avatar_url":             account.AvatarUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute created linkedin connected account %w", err)
	}
	LinkedinItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[linkedin.LinkedinConnectedAccount])
	if err != nil {
		return nil, fmt.Errorf("failed to collect a row from table %w", err)

	}
	return &LinkedinItem, nil
}
