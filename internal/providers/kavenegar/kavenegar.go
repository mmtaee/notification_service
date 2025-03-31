package kavenegar

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"notification/pkg/logger"
	"notification/pkg/phone_parser"
	"os"
)

type Kavenegar struct {
	apiKey            string
	phoneNumberParser phone_parser.ParsedNumberInterface
	validator         *validator.Validate
}

func NewKavenegar() *Kavenegar {
	key := os.Getenv("KAVENEGAR_API_KEY")
	if key == "" {
		logger.Critical("KAVENEGAR_API_KEY environment variable not set")
	}
	return &Kavenegar{
		apiKey:            key,
		phoneNumberParser: phone_parser.NewParser(),
		validator:         validator.New(),
	}
}

func (k *Kavenegar) Process(msg *amqp.Delivery) error {
	switch msg.RoutingKey {
	case "otp.call":
		return k.callOTP(msg)
	case "otp.sms":
		return k.smsOTP(msg)
	case "event.sms":
		return k.eventSms(msg)
	default:
		return fmt.Errorf("invalid routing key: %s for kavenegar provider", msg.RoutingKey)
	}
}

func (k *Kavenegar) smsOTP(msg *amqp.Delivery) error {
	url := fmt.Sprintf("https://api.kavenegar.com/v1/%s/verify/lookup.json", k.apiKey)

	var data SmsOTP
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		return err
	}
	if err = k.validator.Struct(data); err != nil {
		return err
	}

	numObj, err := k.phoneNumberParser.Parse(data.Data.To)
	if err != nil {
		return err
	}
	if !numObj.IsIranNumber() {
		return fmt.Errorf("phone number %s is not iran number", numObj.Number)
	}
	numObj = numObj.IranMasked()

	body, err := json.Marshal(OTPSmsRequest{
		Receptor: numObj.Masked,
		Template: data.Data.Template,
		Token:    data.Data.Code,
	})

	response, err := request(url, "post", body)

	var responseJson OTPSmsResponse
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

func (k *Kavenegar) callOTP(msg *amqp.Delivery) error {
	return errors.New("kavenegar otp call is not yet supported")
}

func (k *Kavenegar) eventSms(msg *amqp.Delivery) error {
	url := fmt.Sprintf("https://api.kavenegar.com/v1/%s/sms/send.json", k.apiKey)
	var data Event
	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		return err
	}
	if err = k.validator.Struct(data); err != nil {
		return err
	}

	numObj, err := k.phoneNumberParser.Parse(data.Data.To)
	if err != nil {
		return err
	}
	if !numObj.IsIranNumber() {
		return fmt.Errorf("phone number %s is not iran number", numObj.Number)
	}
	numObj = numObj.IranMasked()

	body, err := json.Marshal(EventSmsRequest{
		Receptor: numObj.Masked,
		Message:  data.Data.Message,
	})

	response, err := request(url, "post", body)

	var responseJson EventSmsResponse
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
