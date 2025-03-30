package models

type KavenegarOTPSmsRequest struct {
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

type KavenegarOTPSmsResponse struct {
	Return struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"return"`
	Entries []Entry `json:"entries"`
}
