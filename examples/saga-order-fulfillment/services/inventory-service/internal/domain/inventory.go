package domain

import (
	"time"

	"gorm.io/gorm"
)

// ReservationStatus represents inventory reservation state
type ReservationStatus string

const (
	ReservationActive   ReservationStatus = "ACTIVE"
	ReservationReleased ReservationStatus = "RELEASED"
	ReservationFailed   ReservationStatus = "FAILED"
)

// InventoryReservation represents an inventory reservation entity
type InventoryReservation struct {
	ID            string            `gorm:"primaryKey" json:"id"`
	OrderID       string            `gorm:"index;not null" json:"order_id"`
	ProductIDs    string            `gorm:"type:text;not null" json:"product_ids"` // JSON array
	TotalItems    int               `gorm:"not null" json:"total_items"`
	Status        ReservationStatus `gorm:"not null" json:"status"`
	CorrelationID string            `gorm:"index" json:"correlation_id"`
	FailureReason string            `json:"failure_reason,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	ReleasedAt    *time.Time        `json:"released_at,omitempty"`
}

// TableName specifies the table name for GORM
func (InventoryReservation) TableName() string {
	return "inventory_reservations"
}

// BeforeCreate hook to set timestamps
func (r *InventoryReservation) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to update timestamp
func (r *InventoryReservation) BeforeUpdate(tx *gorm.DB) error {
	r.UpdatedAt = time.Now()
	return nil
}

// MarkAsReleased marks reservation as released
func (r *InventoryReservation) MarkAsReleased() {
	r.Status = ReservationReleased
	now := time.Now()
	r.ReleasedAt = &now
	r.UpdatedAt = now
}

// MarkAsFailed marks reservation as failed
func (r *InventoryReservation) MarkAsFailed(reason string) {
	r.Status = ReservationFailed
	r.FailureReason = reason
	r.UpdatedAt = time.Now()
}
