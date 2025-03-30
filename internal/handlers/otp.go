package handlers

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"notification/internal/models"
	"notification/internal/providers/kavenegar"
	"notification/internal/providers/ycloud"
	"notification/pkg/phone_parser"
)

var otpHandlers = map[string]func(*models.OTP, *phone_parser.Parser) bool{
	"otp.sms":  sendSMSOTP,
	"otp.call": nil,
}

func OTPHandler(h *Handler, msg *amqp.Delivery) {
	var data *models.OTP

	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		log.Printf("unmarshal json err: %v", err)
		return
	}
	err = h.validator.Struct(data)
	if err != nil {
		log.Printf("validation error: %s", err)
		return
	}
	numObj, err := h.phoneNumberParser.Parse(data.To)
	if err != nil {
		log.Printf("parse err: %v", err)
	} else {
		if handler, found := otpHandlers[msg.RoutingKey]; found {
			result := handler(data, numObj)
			if result {
				msg.Ack(false)
				return
			}
			log.Println("failed to send, result : ", result)
		}
	}
	log.Println("failed to process : ", string(msg.Body))
	msg.Nack(false, false)
}

func sendSMSOTP(otp *models.OTP, num *phone_parser.Parser) bool {
	switch otp.Provider.Name {
	case "kavenegar":
		if num.IsIranNumber() && otp.Provider.Extra.Template != "" {
			num = num.IranMasked()
			err := kavenegar.SendOTP(num.Masked, otp.Message, otp.Provider.Extra.Template)
			if err != nil {
				log.Println("error sending OTP:", err)
				return false
			}
			return true
		}
		return false
	case "ycloud":
		if num.IsIranNumber() {
			log.Println("failed to send iran numbers")
			return false
		}
		err := ycloud.SendOTP(num.Masked, otp.Message)
		if err != nil {
			log.Println("error sending OTP:", err)
			return false
		}
		return true
	default:
		log.Printf("provider %s not supported\n", otp.Provider.Name)
		return false
	}
}
