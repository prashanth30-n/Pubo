package linkedin

import "github.com/go-playground/validator/v10"

type CreateLinkedinConnectionPayload struct {
	DisplayName string `json:"displayName" validate:"required"`
	AccessToken string `json:"acessToken" validate:"required"`
	PlatformId  *int    `json:"platformId" validate:"required"`
}

func (p *CreateLinkedinConnectionPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}


