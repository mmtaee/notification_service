package kavenegar

type SmsOTP struct {
	Data struct {
		To       string `json:"to" validate:"required"`
		Code     string `json:"code" validate:"required"`
		Template string `json:"template"`
	} `json:"data"`
}

type Event struct {
	Data struct {
		To      string `json:"to" validate:"required"`
		Message string `json:"message" validate:"required"`
	} `json:"data"`
}

type OTPSmsRequest struct {
	Receptor string `json:"receptor"`
	Template string `json:"template"`
	Token    string `json:"token"`
}

type Entry struct {
	MessageID  int    `json:"messageid"`
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statustext"`
	Sender     string `json:"sender"`
	Receptor   string `json:"receptor"`
	Date       int64  `json:"date"`
	Cost       int    `json:"cost"`
}

type OTPSmsResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []Entry `json:"entries"`
}
