package models

import "time"

type YCloudOTPRequest struct {
	To         string `json:"to"`
	Text       string `json:"text"`
	ExternalID string `json:"externalId"`
}

type YCloudOTPResponse struct {
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
