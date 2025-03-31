package ycloud

import "time"

type OTP struct {
	Data struct {
		To       string `json:"to" validate:"required"`
		Code     string `json:"code" validate:"required,min=4,max=8"`
		Brand    string `json:"brand" validate:"required"`
		Channel  string `json:"channel" validate:"omitempty,oneof=sms call email_code"`
		Language string `json:"language" validate:"omitempty,oneof=ar de en es fr id it pt_BR ru tr vi zh_CN zh_HK"`
	} `json:"data"`
}

type Event struct {
	Data struct {
		To      string `json:"to" validate:"required"`
		Message string `json:"message" validate:"required"`
	} `json:"data"`
}

type OTPRequest struct {
	Channel    string `json:"channel"`
	To         string `json:"to"`
	Code       string `json:"text"`
	ExternalID string `json:"externalId"`
	Brand      string `json:"brand"`
	Language   string `json:"language"`
}

type SmsFallback struct {
	Supported         bool   `json:"supported"`
	UnsupportedReason string `json:"unsupportedReason"`
	UnsupportedDetail string `json:"unsupportedDetail"`
}

type OTPResponse struct {
	ID                 string      `json:"id"`
	Status             string      `json:"status"`
	To                 string      `json:"to"`
	Channel            string      `json:"channel"`
	SendTime           time.Time   `json:"sendTime"`
	TotalPrice         float64     `json:"totalPrice"`
	Currency           string      `json:"currency"`
	SmsFallbackEnabled bool        `json:"smsFallbackEnabled"`
	SmsFallback        SmsFallback `json:"smsFallback"`
	ExternalID         string      `json:"externalId"`
}

type EventSmsRequest struct {
	To         string `json:"to"`
	Text       string `json:"text"`
	ExternalID string `json:"externalId"`
}

type EventSmsResponse struct {
	ID             string    `json:"id"`
	To             string    `json:"to"`
	Text           string    `json:"text"`
	SenderID       string    `json:"senderId"`
	RegionCode     string    `json:"regionCode"`
	TotalSegments  int       `json:"totalSegments"`
	TotalPrice     float64   `json:"totalPrice"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	ErrorCode      string    `json:"errorCode"`
	CreateTime     time.Time `json:"createTime"`
	UpdateTime     time.Time `json:"updateTime"`
	ExternalID     string    `json:"externalId"`
	CallbackURL    string    `json:"callbackUrl"`
	BizType        string    `json:"bizType"`
	VerificationID string    `json:"verificationId"`
}
