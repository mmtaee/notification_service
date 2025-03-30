package ycloud

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"notification/internal/models"
)

func SendOTP(to, message string) error {
	url := "https://api.ycloud.com/v2/sms"
	body, _ := json.Marshal(models.YCloudOTPRequest{
		To:         to,
		Text:       message,
		ExternalID: uuid.New().String(),
	})

	response, err := request(url, "post", body)
	if err != nil {
		return err
	}

	var responseJson models.YCloudOTPResponse
	err = json.Unmarshal(response, &responseJson)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	//	TODO: log in db
	return nil
}
