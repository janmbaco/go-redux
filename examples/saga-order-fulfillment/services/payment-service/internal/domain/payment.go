package domain

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents payment state
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentRefunded  PaymentStatus = "REFUNDED"
	PaymentFailed    PaymentStatus = "FAILED"
)

// Payment represents a payment entity
type Payment struct {
	ID            string        `gorm:"primaryKey" json:"id"`
	OrderID       string        `gorm:"index;not null" json:"order_id"`
	CustomerID    string        `gorm:"index;not null" json:"customer_id"`
	Amount        float64       `gorm:"not null" json:"amount"`
	Status        PaymentStatus `gorm:"not null" json:"status"`
	CorrelationID string        `gorm:"index" json:"correlation_id"`
	FailureReason string        `json:"failure_reason,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	RefundedAt    *time.Time    `json:"refunded_at,omitempty"`
}

// TableName specifies the table name for GORM
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate hook to set timestamps
func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to update timestamp
func (p *Payment) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// MarkAsCompleted marks payment as completed
func (p *Payment) MarkAsCompleted() {
	p.Status = PaymentCompleted
	p.UpdatedAt = time.Now()
}

// MarkAsRefunded marks payment as refunded
func (p *Payment) MarkAsRefunded() {
	p.Status = PaymentRefunded
	now := time.Now()
	p.RefundedAt = &now
	p.UpdatedAt = now
}

// MarkAsFailed marks payment as failed
func (p *Payment) MarkAsFailed(reason string) {
	p.Status = PaymentFailed
	p.FailureReason = reason
	p.UpdatedAt = time.Now()
}
