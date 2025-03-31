package ycloud

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification/pkg/logger"
	"notification/pkg/phone_parser"
	"os"
)

type Ycloud struct {
	apiKey            string
	phoneNumberParser phone_parser.ParsedNumberInterface
	validator         *validator.Validate
}

func NewYcloud() *Ycloud {
	key := os.Getenv("YCLOUD_API_KEY")
	if key == "" {
		logger.Critical("YCloud_API_KEY environment variable not set")
	}
	return &Ycloud{
		apiKey:            key,
		phoneNumberParser: phone_parser.NewParser(),
		validator:         validator.New(),
	}
}

func (y *Ycloud) Process(msg *amqp.Delivery) error {
	switch msg.RoutingKey {
	case "otp.call":
		return y.OTP(msg)
	case "otp.sms":
		return y.OTP(msg)
	case "event.sms":
		//return y.eventSms(msg)
		return y.eventSms(msg)
	default:
		return fmt.Errorf("invalid routing key: %s for kavenegar provider", msg.RoutingKey)
	}
}

func (y *Ycloud) OTP(msg *amqp.Delivery) error {
	var data OTP

	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		return err
	}
	if data.Data.Channel == "" {
		data.Data.Channel = "sms"
	}
	if data.Data.Language == "" {
		data.Data.Language = "en"
	}
	if err = y.validator.Struct(data); err != nil {
		return err
	}

	numObj, err := y.phoneNumberParser.Parse(data.Data.To)
	if err != nil {
		return err
	}

	if numObj.IsIranNumber() {
		return fmt.Errorf("iran number %s does not allowd to send sms with ycloud provider", numObj.Masked)
	}

	url := "https://api.ycloud.com/v2/verify/verifications"
	body, _ := json.Marshal(OTPRequest{
		Channel:    data.Data.Channel,
		To:         numObj.Masked,
		Code:       data.Data.Code,
		Brand:      data.Data.Brand,
		ExternalID: uuid.New().String(),
		Language:   data.Data.Language,
	})

	response, err := request(url, "post", body, y.apiKey)
	if err != nil {
		return err
	}

	var responseJson OTPResponse
	err = json.Unmarshal(response, &responseJson)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	//	TODO: log in db
	return nil
}
func (y *Ycloud) eventSms(msg *amqp.Delivery) error {
	var data Event
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		return err
	}
	if err = y.validator.Struct(data); err != nil {
		return err
	}
	numObj, err := y.phoneNumberParser.Parse(data.Data.To)
	if err != nil {
		return err
	}

	if numObj.IsIranNumber() {
		return fmt.Errorf("iran number %s does not allowd to send sms with ycloud provider", numObj.Masked)
	}

	url := "https://api.ycloud.com/v2/sms"

	body, _ := json.Marshal(EventSmsRequest{
		To:         numObj.Masked,
		Text:       data.Data.Message,
		ExternalID: uuid.New().String(),
	})
	response, err := request(url, "post", body, y.apiKey)
	if err != nil {
		return err
	}

	var responseJson EventSmsResponse
	err = json.Unmarshal(response, &responseJson)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	//	TODO: log in db
	return nil
}
