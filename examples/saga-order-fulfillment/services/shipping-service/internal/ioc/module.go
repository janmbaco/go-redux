package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/shipping-service/internal/domain"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/shipping-service/internal/handlers"
	"gorm.io/gorm"
)

// ShippingModule implements Module for shipping service
type ShippingModule struct{}

// NewShippingModule creates a new shipping module
func NewShippingModule() *ShippingModule {
	return &ShippingModule{}
}

// RegisterServices registers shipping service dependencies
func (m *ShippingModule) RegisterServices(register dependencyinjection.Register) error {
	// Register ShippingHandler
	register.AsSingleton(
		new(*handlers.ShippingHandler),
		func(db *gorm.DB, logger logs.Logger) (*handlers.ShippingHandler, error) {
			// Auto-migrate
			if err := db.AutoMigrate(&domain.Shipment{}); err != nil {
				return nil, err
			}
			// Create typed data access
			da := dataaccess.NewTypedDataAccess[domain.Shipment](db)
			return handlers.NewShippingHandler(da, logger), nil
		},
		map[int]string{0: "db", 1: "logger"},
	)

	return nil
}
