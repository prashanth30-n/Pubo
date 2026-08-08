package signup

import "github.com/prashanth30-n/Pubo/internal/model"

type SignUp struct {
	model.Base
	ClerkUserId *string `json:"userID" db:"clerk_user_id"`
	Email       *string `json:"email" db:"email"`
	DisplayName string  `json:"displayName" db:"display_name"`
	AvatarUrl   string  `json:"avatarUrl" db:"avatar_url"`
	isActive    bool    `json"isActive" db:"is_active"`
}
