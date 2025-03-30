package models

type OTP struct {
	To       string `json:"to" validate:"required"`
	Message  string `json:"message" validate:"required"`
	Provider struct {
		Name  string `json:"name" validate:"required,min=1"`
		Extra struct {
			Template string `json:"template"`
		} `json:"extra" validate:"omitempty"`
	} `json:"provider" validate:"required"`
}

type SMS struct {
	To       string `json:"to" validate:"required"`
	Message  string `json:"message" validate:"required"`
	Provider string `json:"provider" validate:"required"`
}
