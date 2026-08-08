package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/prashanth30-n/Pubo/internal/server"
)

type ConnectedAccount struct {
	ID          string
	Platform    string
	DID         string
	AccessToken string
	PDSURL      string
}
type PostRepository struct{ server *server.Server }

type SavedPost struct {
	ID        string      `json:"id"`
	Content   string      `json:"content"`
	Status    string      `json:"status"`
	MediaIDs  []string    `json:"mediaIds"`
	Platforms []string    `json:"platforms"`
	CreatedAt string      `json:"createdAt"`
	Media     []PostMedia `json:"media"`
}
type PostMedia struct {
	ID          string `json:"id"`
	StorageURL  string `json:"storageURL"`
	ContentType string `json:"contentType"`
}
type PostPage struct {
	Data   []SavedPost `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}

func NewPostRepository(s *server.Server) *PostRepository { return &PostRepository{server: s} }

func (r *PostRepository) GetActiveAccount(ctx context.Context, userID, platform string) (*ConnectedAccount, error) {
	const query = `SELECT ca.id, p.code, COALESCE(ca.did,''), COALESCE(ca.access_token_encrypted,''), COALESCE(ca.pds_url,'')
		FROM connected_accounts ca JOIN platforms p ON p.id=ca.platform_id
		WHERE ca.user_id=@user_id AND p.code=@platform AND ca.is_active=true LIMIT 1`
	var account ConnectedAccount
	err := r.server.DB.Pool.QueryRow(ctx, query, pgx.NamedArgs{"user_id": userID, "platform": platform}).Scan(&account.ID, &account.Platform, &account.DID, &account.AccessToken, &account.PDSURL)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no active %s account is connected", platform)
	}
	if err != nil {
		return nil, err
	}
	if account.AccessToken == "" {
		return nil, fmt.Errorf("%s account has no access token", platform)
	}
	return &account, nil
}

// RecordPlatformStatus creates an audit post/target only after the provider call.
// It keeps a failed provider from looking like a successfully published post.
func (r *PostRepository) RecordPlatformStatus(ctx context.Context, userID, platform, content, status, externalPostID, failureReason string) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var postID string
	err = tx.QueryRow(ctx, `INSERT INTO posts(author_user_id,content,status,posted_at,failure_reason) SELECT id,@content,@status::post_status,CASE WHEN @status::post_status='posted'::post_status THEN NOW() ELSE NULL END,@failure FROM pubo_users WHERE clerk_user_id=@user_id RETURNING id`, pgx.NamedArgs{"content": content, "status": status, "failure": failureReason, "user_id": userID}).Scan(&postID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO post_platform_status(post_id,connected_account_id,platform_code,external_post_id,status,failure_reason,published_at) SELECT @post_id,ca.id,@platform,@external,@status::post_status,@failure,CASE WHEN @status::post_status='posted'::post_status THEN NOW() ELSE NULL END FROM connected_accounts ca JOIN platforms p ON p.id=ca.platform_id WHERE ca.user_id=@user_id AND p.code=@platform`, pgx.NamedArgs{"post_id": postID, "platform": platform, "external": externalPostID, "status": status, "failure": failureReason, "user_id": userID})
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *PostRepository) CreateDraft(ctx context.Context, userID, content string, mediaIDs []string) (*SavedPost, error) {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	post := &SavedPost{Content: content, Status: "draft", MediaIDs: mediaIDs}
	err = tx.QueryRow(ctx, `INSERT INTO posts(author_user_id, content, status) SELECT id,@content,'draft'::post_status FROM pubo_users WHERE clerk_user_id=@user_id RETURNING id,created_at::text`, pgx.NamedArgs{"content": content, "user_id": userID}).Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	for position, mediaID := range mediaIDs {
		_, err = tx.Exec(ctx, `INSERT INTO post_media_links(post_id,media_id,position) VALUES(@post_id,@media_id,@position)`, pgx.NamedArgs{"post_id": post.ID, "media_id": mediaID, "position": position})
		if err != nil {
			return nil, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) ListPosts(ctx context.Context, userID, status string, limit, offset int) (*PostPage, error) {
	query := `SELECT p.id::text,p.content,p.status::text,COALESCE(array_agg(DISTINCT pml.media_id::text) FILTER (WHERE pml.media_id IS NOT NULL),'{}'),COALESCE(array_agg(DISTINCT pps.platform_code) FILTER (WHERE pps.platform_code IS NOT NULL),'{}'),p.created_at::text FROM posts p JOIN pubo_users u ON u.id=p.author_user_id LEFT JOIN post_media_links pml ON pml.post_id=p.id LEFT JOIN post_platform_status pps ON pps.post_id=p.id WHERE u.clerk_user_id=@user_id AND p.status=COALESCE(NULLIF(@status,'')::post_status,p.status) GROUP BY p.id ORDER BY p.created_at DESC`
	query += ` LIMIT @limit OFFSET @offset`
	rows, err := r.server.DB.Pool.Query(ctx, query, pgx.NamedArgs{"user_id": userID, "status": status, "limit": limit, "offset": offset})
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SavedPost{}
	for rows.Next() {
		var post SavedPost
		if err := rows.Scan(&post.ID, &post.Content, &post.Status, &post.MediaIDs, &post.Platforms, &post.CreatedAt); err != nil {
			return nil, err
		}
		mediaRows, err := r.server.DB.Pool.Query(ctx, `SELECT pm.id::text,pm.storage_url,pm.content_type FROM post_media_links pml JOIN post_media pm ON pm.id=pml.media_id WHERE pml.post_id=@post_id ORDER BY pml.position`, pgx.NamedArgs{"post_id": post.ID})
		if err != nil {
			return nil, err
		}
		for mediaRows.Next() {
			var media PostMedia
			if err := mediaRows.Scan(&media.ID, &media.StorageURL, &media.ContentType); err != nil {
				mediaRows.Close()
				return nil, err
			}
			post.Media = append(post.Media, media)
		}
		mediaRows.Close()
		out = append(out, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	var total int
	if err := r.server.DB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM posts p JOIN pubo_users u ON u.id=p.author_user_id WHERE u.clerk_user_id=@user_id AND p.status=COALESCE(NULLIF(@status,'')::post_status,p.status)`, pgx.NamedArgs{"user_id": userID, "status": status}).Scan(&total); err != nil {
		return nil, err
	}
	return &PostPage{Data: out, Limit: limit, Offset: offset, Total: total}, nil
}
