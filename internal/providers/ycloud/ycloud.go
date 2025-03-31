package ycloud

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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
		log.Fatal("YCloud_API_KEY environment variable not set")
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
		return nil
		//return y.callOTP(msg)
	case "otp.sms":
		return y.smsOTP(msg)
	case "event.sms":
		//return y.eventSms(msg)
		return nil
	default:
		return fmt.Errorf("invalid routing key: %s for kavenegar provider", msg.RoutingKey)
	}
}

func (y *Ycloud) smsOTP(msg *amqp.Delivery) error {
	var data SmsOTP

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
	body, _ := json.Marshal(SMSOTPRequest{
		To:         numObj.Masked,
		Text:       data.Data.Message,
		ExternalID: uuid.New().String(),
	})

	response, err := request(url, "post", body, y.apiKey)
	if err != nil {
		return err
	}

	var responseJson SMSOTPResponse
	err = json.Unmarshal(response, &responseJson)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return err
	}
	//	TODO: log in db
	return nil
}
