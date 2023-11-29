package models

import "time"

type Payment struct {
	ID            uint      `gorm:"primaryKey"`
	Amount        float64   `gorm:"not null"`
	PaymentType   int64     `gorm:"not null"`
	PaymentStatus int64     `gorm:"not null"`
	OrderID       uint      `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	// Payment belongs to an Order
	Order Order `gorm:"foreignKey:OrderID"`
}

type AdvancePaymentRequest struct {
	OrderId uint `json:"order_id"`
}

type PaymentRequest struct {
	MerchantId    string `json:"MerchantId"`
	TerminalId    string `json:"TerminalId"`
	ReturnUrl     string `json:"ReturnUrl"`
	LocalDateTime string `json:"LocalDateTime"`
	SignData      string `json:"SignData"`
	OrderId       uint  `json:"OrderId"`
	Amount        int    `json:"Amount"`
}

type MelliPaymentResponse struct {
	ResCode     int    `json:"ResCode"`
	Token       string `json:"Token"`
	Description string `json:"Description"`
	URL         string `json:"url"`
}

type VerifyPaymentRequest struct {
	SignData string `json:"SignData"`
	Token    string `json:"Token"`
}

type MelliVerifyPaymentResponse struct {
	ResCode       int    `json:"ResCode"`
	Amount        int    `json:"Amount"`
	Description   string `json:"Description"`
	RetrivalRefNo string `json:"RetrivalRefNo"`
	SystemTraceNo string `json:"SystemTraceNo"`
	OrderId       int64  `json:"OrderId"`
}

const (
	AdvancePayment = iota
	FinalPayment
)

const (
	PaymentCompleted = iota
	PaymentPending
	PaymentFailed
)
