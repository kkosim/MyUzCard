package models

import "time"

// RequestLogin структура для хранения логина и пароля
type RequestLogin struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

// Response структура для хранения ответа от Узбекистана
type Response struct {
	Results *Results `json:"result"`
	Error   *Error   `json:"error"`
	Success bool     `json:"success"`
}

// Error структура для хранения вывода ошибок
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Results структура для хранения всех Result
type Results struct {
	*ResultLogin
	*ResultCheck
	*ResultPay
	*ResultGetStatus
	*ResultReverse
	*ResultP2pinfowrap
}

// ResultLogin структура для хранения содержания ответа на /login
type ResultLogin struct {
	Token      string `json:"token"`
	UserID     int64  `json:"userId"`
	ExpireDate string `json:"expireDate"`
	RoleID     int64  `json:"roleId"`
}

// CurrencyInfo структура для хранения информации о валюте
type CurrencyInfo struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

// ResultCheck структура для хранения содержания ответа на /check
type ResultCheck struct {
	TransactionNumber string        `json:"transactionNumber"`
	State             int           `json:"state"`
	CreatedDate       string        `json:"createdDate"`
	Amount            int           `json:"amount"`
	CommissionAmount  int           `json:"commissionAmount"`
	TotalAmount       int           `json:"totalAmount"`
	CardHolder        string        `json:"cardHolder"`
	BankName          string        `json:"bankName"`
	CurrencyInfo      *CurrencyInfo `json:"currencyInfo"`
}

// SenderDocument структура для хранения информации о документе отправителя
type SenderDocument struct {
	Nationality string `json:"nationality"`
	DocType     string `json:"docType"`
	DocSerial   string `json:"docSerial"`
	DocNumber   string `json:"docNumber"`
	ValidTo     string `json:"validTo"`
	BirthDate   string `json:"birthDate"`
}

// Sender структура для хранения информации об отправителе
type Sender struct {
	CardNumber     string          `json:"cardNumber"`
	CardIssuer     string          `json:"cardIssuer"`
	Country        string          `json:"country"`
	SenderDocument *SenderDocument `json:"senderDocument"`
}

// RequestCheck структура для хранения запроса на /check
type RequestCheck struct {
	Amount                   float64 `json:"amount"`
	CardNumber               string  `json:"cardNumber"`
	Currency                 string  `json:"currency"`
	PartnerClientId          string  `json:"partnerClientId"`
	PartnerTransactionNumber string  `json:"partnerTransactionNumber"`
	SenderFirstName          string  `json:"senderFirstName"`
	SenderMiddleName         *string `json:"senderMiddleName"`
	SenderLastName           string  `json:"senderLastName"`
	SenderAddress            *string `json:"senderAddress"`
	Sender                   *Sender `json:"sender"`
}

// ResultPay структура для хранения содержания ответа на /pay
type ResultPay struct {
	TransactionNumber string  `json:"transactionNumber"`
	State             int     `json:"state"`
	CreatedDate       string  `json:"createdDate"`
	PaymentDate       string  `json:"paymentDate"`
	Amount            float64 `json:"amount"`
	TotalAmount       float64 `json:"totalAmount"`
	CommissionAmount  float64 `json:"commissionAmount"`
}

// TransactionNumber структура для хранения номера транзакции
type TransactionNumber struct {
	TransactionNumber string `json:"transactionNumber"`
}

// ResultGetStatus структура для хранения содержания ответа на /get-status
type ResultGetStatus struct {
	State            int     `json:"state"`
	CreatedDate      string  `json:"createdDate"`
	PaymentDate      string  `json:"paymentDate"`
	Amount           float64 `json:"amount"`
	TotalAmount      float64 `json:"totalAmount"`
	CommissionAmount float64 `json:"commissionAmount"`
}

// ResultReverse структура для хранения содержания ответа на /reverse
type ResultReverse struct {
	State            int       `json:"state"`
	CreatedDate      string    `json:"createdDate"`
	PaymentDate      string    `json:"paymentDate"`
	Amount           float64   `json:"amount"`
	TotalAmount      float64   `json:"totalAmount"`
	CommissionAmount float64   `json:"commissionAmount"`
	ReverseDate      time.Time `json:"reverseDate"`
}

// CardNumber структура для хранения номера карты
type CardNumber struct {
	CardNumber string `json:"cardNumber"`
}

// ResultP2pinfowrap структура для хранения информации о карте
type ResultP2pinfowrap struct {
	CardHolder string `json:"cardHolder"`
	BankName   string `json:"bankName"`
}
