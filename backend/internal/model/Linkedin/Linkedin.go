package linkedin

import (
	"github.com/PatibandlaVenkat/Pubo/internal/model"
	"time"
)

type LinkedinConnectedAccount struct {
	model.Base
	UserID           string  `json:"userID" db:"user_id"`
	PlatformId       int    `json:"PlatformId" db:"platform_id"`
	Handle           *string `json:"Handle" db:"handle"`
	DisplayName      string `json:"DisplayName" db:"display_name"`
	AvatarUrl        *string   `json:"avatarUrl" db:"avatar_url"`
	Password         *string `json:"Password"`
	AccessToken  string `json:"acessToken" db:"access_token_encrypted"`
	RefreshToken  *string  `json:"refreshToken" db:"refresh_token_encrypted"`
	TokenExpiry *time.Time `json:"tokenExpiresAt" db:"token_expires_at"`
	DID  string `json:"did" db:"did"`
	PDSURL string `json:"pds_url" db:"pds_url"`
	Scopes    *[]string `json:"scopes" db:"scopes"`
	Is_active      bool `json:"isActive" db:"is_active"`
	Last_synced_at   *string `json:"lastSyncedAt" db:"last_synced_at"`
}


