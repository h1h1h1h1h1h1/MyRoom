package models

import (
	"time"
)

type CustomerNumber struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	CustomerNumber   string    `json:"customer_number"`
	MeterNumber      string    `json:"meter_number"`
	Address          string    `json:"address"`
	VoltageLevel     string    `json:"voltage_level,omitempty"`
	ContractCapacity float64   `json:"contract_capacity,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ElectricityBalance struct {
	ID               int       `json:"id"`
	CustomerNumberID int       `json:"customer_number_id"`
	Balance          float64   `json:"balance"`
	LastUpdated      time.Time `json:"last_updated"`
}

type MeterReading struct {
	ID               int       `json:"id"`
	CustomerNumberID int       `json:"customer_number_id"`
	ReadingDate      time.Time `json:"reading_date"`
	CurrentReading   float64   `json:"current_reading"`
	PreviousReading  float64   `json:"previous_reading,omitempty"`
	Consumption      float64   `json:"consumption"`
	ReadingType      string    `json:"reading_type"`
	ReaderID         int       `json:"reader_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type ElectricityBill struct {
	ID               int       `json:"id"`
	CustomerNumberID int       `json:"customer_number_id"`
	BillingMonth     time.Time `json:"billing_month"`
	Consumption      float64   `json:"consumption"`
	UnitPrice        float64   `json:"unit_price"`
	Amount           float64   `json:"amount"`
	DueDate          time.Time `json:"due_date"`
	PaidAmount       float64   `json:"paid_amount"`
	PaymentStatus    string    `json:"payment_status"`
	PaidDate         time.Time `json:"paid_date,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type PurchaseRecord struct {
	ID               int       `json:"id"`
	CustomerNumberID int       `json:"customer_number_id"`
	PurchaseAmount   float64   `json:"purchase_amount"`
	PaymentMethod    string    `json:"payment_method"`
	TransactionID    string    `json:"transaction_id,omitempty"`
	PurchaseDate     time.Time `json:"purchase_date"`
	Status           string    `json:"status"`
	Notes            string    `json:"notes,omitempty"`
}

// Notification struct moved to common.go
// type Notification struct { ... }

type InformationPost struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	PostType    string    `json:"post_type"`
	PublishDate time.Time `json:"publish_date"`
	ExpiryDate  time.Time `json:"expiry_date,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   int       `json:"created_by,omitempty"`
}

type ElectricityApplication struct {
	ID                int       `json:"id"`
	UserID            int       `json:"user_id"`
	ApplicationType   string    `json:"application_type"`
	CustomerNumberID  int       `json:"customer_number_id,omitempty"`
	ApplicationDate   time.Time `json:"application_date"`
	Status            string    `json:"status"`
	Description       string    `json:"description,omitempty"`
	RequiredDocuments string    `json:"required_documents,omitempty"`
	ProcessingNotes   string    `json:"processing_notes,omitempty"`
	CompletedDate     time.Time `json:"completed_date,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type ServicePoint struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Address       string  `json:"address"`
	Phone         string  `json:"phone,omitempty"`
	BusinessHours string  `json:"business_hours,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	IsActive      bool    `json:"is_active"`
}

type ElectricityPrice struct {
	ID             int       `json:"id"`
	VoltageLevel   string    `json:"voltage_level"`
	Tier           int       `json:"tier"`
	MinConsumption float64   `json:"min_consumption,omitempty"`
	MaxConsumption float64   `json:"max_consumption,omitempty"`
	UnitPrice      float64   `json:"unit_price"`
	EffectiveDate  time.Time `json:"effective_date"`
	ExpiryDate     time.Time `json:"expiry_date,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// 请求和响应结构
type BindCustomerRequest struct {
	CustomerNumber string `json:"customer_number" binding:"required"`
	MeterNumber    string `json:"meter_number" binding:"required"`
	Address        string `json:"address" binding:"required"`
}

type PurchaseRequest struct {
	CustomerNumberID int     `json:"customer_number_id" binding:"required"`
	PurchaseAmount   float64 `json:"purchase_amount" binding:"required,gt=0"`
	PaymentMethod    string  `json:"payment_method" binding:"required"`
}

type MeterReadingRequest struct {
	CustomerNumberID int       `json:"customer_number_id" binding:"required"`
	ReadingDate      time.Time `json:"reading_date" binding:"required"`
	CurrentReading   float64   `json:"current_reading" binding:"required"`
	ReadingType      string    `json:"reading_type"`
}

type BillPaymentRequest struct {
	BillID     int     `json:"bill_id" binding:"required"`
	PaidAmount float64 `json:"paid_amount" binding:"required,gt=0"`
}

// ApplicationRequest struct moved to common.go
// type ApplicationRequest struct { ... }

type ElectricityOverview struct {
	CustomerNumber     string    `json:"customer_number"`
	Address            string    `json:"address"`
	Balance            float64   `json:"balance"`
	LastReadingDate    time.Time `json:"last_reading_date,omitempty"`
	LastReading        float64   `json:"last_reading,omitempty"`
	LatestUnpaidAmount float64   `json:"latest_unpaid_amount,omitempty"`
}
