package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/prashanth30-n/Pubo/internal/model/signup"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type SignupRepository struct {
	server *server.Server
}

func NewSignupRepository(server *server.Server) *SignupRepository {
	return &SignupRepository{
		server: server,
	}
}
func (r *SignupRepository) SignUp(ctx context.Context, payload *signup.SignUpPayload, userID string) (*signup.SignUp, error) {
	stmt := `INSERT INTO
	pubo_users(
	clerk_user_id,
	email
	) VALUES(
	 @clerk_user_id,
	 @email
	 )
	 RETURNING
	 *
	 `
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"clerk_user_id": userID,
		"email":         payload.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute signup query: %w", err)

	}
	SignupItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[signup.SignUp])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from database:%w", err)
	}
	return &SignupItem, nil

}
