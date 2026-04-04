package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/payment-service/internal/domain"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/payment-service/internal/handlers"
	"gorm.io/gorm"
)

// PaymentModule implements Module for payment service
type PaymentModule struct{}

// NewPaymentModule creates a new payment module
func NewPaymentModule() *PaymentModule {
	return &PaymentModule{}
}

// RegisterServices registers payment service dependencies
func (m *PaymentModule) RegisterServices(register dependencyinjection.Register) error {
	// Register PaymentHandler
	register.AsSingleton(
		new(*handlers.PaymentHandler),
		func(db *gorm.DB, logger logs.Logger) (*handlers.PaymentHandler, error) {
			// Auto-migrate
			if err := db.AutoMigrate(&domain.Payment{}); err != nil {
				return nil, err
			}
			// Create typed data access
			da := dataaccess.NewTypedDataAccess[domain.Payment](db)
			return handlers.NewPaymentHandler(da, logger), nil
		},
		map[int]string{0: "db", 1: "logger"},
	)

	return nil
}
