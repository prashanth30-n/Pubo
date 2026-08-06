package repository

import (
	"context"
	"errors"

	"github.com/PatibandlaVenkat/Pubo/internal/model/media"
	"github.com/PatibandlaVenkat/Pubo/internal/server"
	"github.com/jackc/pgx/v5"
)

type MediaRepo interface {
	create(ctx context.Context, asset media.Asset) error
	GetByID(ctx context.Context,id string)(*media.Asset,error)
	ListByOwner(ctx context.Context,ownerUserID string,limit int ,offset int)
	Delete(ctx context.Context,id string) error
	
}
type MediaRepository struct{                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        
server *server.Server
}
func NewMediaRepository(server *server.Server) *MediaRepository{
	return  &MediaRepository{
		server: server,
	}
}
func(r*MediaRepository) Create(ctx context.Context,asset media.Asset) error{
	stmt:=`INSERT INTO post_media (id,owner_user_id,original_name,blob_name,storage_url,container_name,content_type,size_bytes,etag,created_at,updated_at)
	VALUES (
	@id,
	@owner_user_id,
	@original_name,
	@blob_name,
	@storage_url,
	@container_name,
	@content_type,
	@size_bytes,
	@etag,
	NOW(),
	NOW()
	)
	
	`
	_,err:=r.server.DB.Pool.Exec(ctx,stmt,pgx.NamedArgs{
		"id":asset.ID,
		"owner_user_id":asset.OwnerUserID,
		"original_name":asset.OriginalName,
		"blob_name":asset.BlobName,
		"storage_url":asset.StorageURL,
		"container_name":asset.ContainerName,
		"content_type":asset.ContentType,
		"size_bytes":asset.SizeBytes,
		"etag":asset.ETag,
	})
	return err
	
}
func(r*MediaRepository) GetByID(ctx context.Context,id string)(*media.Asset,error){
	stmt:=`SELECT id,owner_user_id,original_name,blob_name,storage_url,container_name,content_type,size_bytes,etag,created_at,updated_at FROM post_media
	WHERE id=@id`
	rows,err:=r.server.DB.Pool.Query(ctx,stmt,pgx.NamedArgs{
		"id":id,
	})
	if err!=nil{
		if errors.Is(err,pgx.ErrNoRows){
			return nil,nil
		}
		return nil,err
	}
	asset,err:=pgx.CollectExactlyOneRow(rows,pgx.RowToStructByName[media.Asset])
	if err!=nil{
		return nil,nil
	}
	return &asset,nil
}

func (r *MediaRepository) ListByOwner(ctx context.Context,ownerUserID string, limit int, offset int) ([]media.Asset, error) {
    rows, err := r.server.DB.Pool.Query(ctx, `
        SELECT id, owner_user_id, original_name, blob_name, storage_url, container_name, content_type, size_bytes, etag, created_at, updated_at
        FROM post_media
        WHERE owner_user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `, ownerUserID, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    out := make([]media.Asset, 0)
    for rows.Next() {
        var asset media.Asset
        if err := rows.Scan(
            &asset.ID, &asset.OwnerUserID, &asset.OriginalName, &asset.BlobName,
            &asset.StorageURL, &asset.ContainerName, &asset.ContentType,
            &asset.SizeBytes, &asset.ETag, &asset.CreatedAt, &asset.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        out = append(out, asset)
    }
    return out, rows.Err()
}
func(r *MediaRepository) Delete(ctx context.Context,id string) error{
	_,err:=r.server.DB.Pool.Exec(ctx,`
	DELETE FROM post_media
	WHERE id=@id`,pgx.NamedArgs{
		"id":id,
	})
	return err
}

