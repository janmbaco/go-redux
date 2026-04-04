package ioc

import (
	"github.com/janmbaco/go-infrastructure/v2/dependencyinjection"
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/inventory-service/internal/domain"
	"github.com/janmbaco/go-redux-examples/saga-order-fulfillment/services/inventory-service/internal/handlers"
	"gorm.io/gorm"
)

// InventoryModule implements Module for inventory service
type InventoryModule struct{}

// NewInventoryModule creates a new inventory module
func NewInventoryModule() *InventoryModule {
	return &InventoryModule{}
}

// RegisterServices registers inventory service dependencies
func (m *InventoryModule) RegisterServices(register dependencyinjection.Register) error {
	// Register InventoryHandler
	register.AsSingleton(
		new(*handlers.InventoryHandler),
		func(db *gorm.DB, logger logs.Logger) (*handlers.InventoryHandler, error) {
			// Auto-migrate
			if err := db.AutoMigrate(&domain.InventoryReservation{}); err != nil {
				return nil, err
			}
			// Create typed data access
			da := dataaccess.NewTypedDataAccess[domain.InventoryReservation](db)
			return handlers.NewInventoryHandler(da, logger), nil
		},
		map[int]string{0: "db", 1: "logger"},
	)

	return nil
}
