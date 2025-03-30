package kavenegar

import (
	"encoding/json"
	"errors"
	"fmt"
	"notification/internal/models"
	"notification/pkg/configs"
)

func SendOTP(to, message, template string) error {
	apiKey := configs.KavenegarApiKey()
	url := fmt.Sprintf("https://api.kavenegar.com/v1/%s/verify/lookup.json", apiKey)
	body, err := json.Marshal(models.KavenegarOTPSmsRequest{
		Receptor: to,
		Template: template,
		Token:    message,
	})

	response, err := request(url, "post", body)

	var responseJson models.KavenegarOTPSmsResponse
	err = json.Unmarshal(response, &responseJson)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	if responseJson.Return.Status != 200 {
		return errors.New(responseJson.Return.Message)
	}
	if responseJson.Entries != nil {
		//	TODO: log in db
	}
	return nil
}
