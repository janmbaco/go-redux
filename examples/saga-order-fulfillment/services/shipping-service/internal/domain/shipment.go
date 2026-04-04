package domain

import (
	"time"

	"gorm.io/gorm"
)

// ShipmentStatus represents shipment state
type ShipmentStatus string

const (
	ShipmentCreated   ShipmentStatus = "CREATED"
	ShipmentCancelled ShipmentStatus = "CANCELLED"
	ShipmentFailed    ShipmentStatus = "FAILED"
)

// Shipment represents a shipment entity
type Shipment struct {
	ID            string         `gorm:"primaryKey" json:"id"`
	OrderID       string         `gorm:"index;not null" json:"order_id"`
	Address       string         `gorm:"type:text;not null" json:"address"`
	Status        ShipmentStatus `gorm:"not null" json:"status"`
	CorrelationID string         `gorm:"index" json:"correlation_id"`
	TrackingCode  string         `json:"tracking_code,omitempty"`
	FailureReason string         `json:"failure_reason,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	CancelledAt   *time.Time     `json:"cancelled_at,omitempty"`
}

// TableName specifies the table name for GORM
func (Shipment) TableName() string {
	return "shipments"
}

// BeforeCreate hook to set timestamps
func (s *Shipment) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to update timestamp
func (s *Shipment) BeforeUpdate(tx *gorm.DB) error {
	s.UpdatedAt = time.Now()
	return nil
}

// MarkAsCancelled marks shipment as cancelled
func (s *Shipment) MarkAsCancelled() {
	s.Status = ShipmentCancelled
	now := time.Now()
	s.CancelledAt = &now
	s.UpdatedAt = now
}

// MarkAsFailed marks shipment as failed
func (s *Shipment) MarkAsFailed(reason string) {
	s.Status = ShipmentFailed
	s.FailureReason = reason
	s.UpdatedAt = time.Now()
}
