package handlers

import (
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"notification/pkg/phone_parser"
	"slices"
	"strings"
)

type Handler struct {
	validator         *validator.Validate
	phoneNumberParser phone_parser.ParsedNumberInterface
}

var routingKeys = []string{"otp.sms", "otp.call", "event.sms"}

func NewHandler() *Handler {
	return &Handler{
		validator:         validator.New(),
		phoneNumberParser: phone_parser.NewParser(),
	}
}

func (h *Handler) Send(msg *amqp.Delivery) {
	log.Printf("message: %s", msg.Body)
	if slices.Contains(GetRoutingKeys(), msg.RoutingKey) {
		if strings.HasPrefix(msg.RoutingKey, "otp.") {
			OTPHandler(h, msg)
		} else if strings.HasPrefix(msg.RoutingKey, "event.sms") {
			EventSMSHandler(h, msg)
		} else if strings.HasPrefix(msg.RoutingKey, "event.push") {
			EventPushHandler(h, msg)
		} else {
			log.Println("routing key not found: ", msg.RoutingKey)
		}
		return
	}
	log.Printf("not found routing key: %s", msg.RoutingKey)
	msg.Nack(false, false)
	return
}

func GetRoutingKeys() []string {
	return routingKeys
}
