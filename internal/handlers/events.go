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

func EventSMSHandler(h *Handler, msg *amqp.Delivery) {
	var data *models.SMS
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
		sendEventSMS(data, numObj)
	}
	log.Printf("not found handler: %s", msg.RoutingKey)

}

func EventPushHandler(h *Handler, msg *amqp.Delivery) {

}

func sendEventSMS(sms *models.SMS, num *phone_parser.Parser) bool {
	switch sms.Provider {
	case "kavenegar":
		if num.IsIranNumber() {
			num = num.IranMasked()
			err := kavenegar.SendSMS(num.Masked, sms.Message)
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
		err := ycloud.SendSMS(num.Masked, sms.Message)
		if err != nil {
			log.Println("error sending OTP:", err)
			return false
		}
		return true
	default:
		log.Printf("provider %s not supported\n", sms.Provider)
		return false
	}
}

func sendEventPush() bool {
	return true
}
