package app

import (
	"os"
)

type DatabaseConfig struct {
	User         string
	Password     string
	Address      string
	DatabaseName string
}

type Config struct {
	InventoryDatabaseConfig DatabaseConfig
	BackendUrl              string
	OrderServiceUrl         string
}

// Required Configs:
// - DB_ADDRESS
// - DB_USER
// - DB_PASSWORD
// - DB_NAME
// - BACKEND_URL
// - ORDER_SERVICE_URL
func LoadConfig() (*Config, error) {
	cfg := &Config{
		InventoryDatabaseConfig: DatabaseConfig{
			User:     "user",
			Password: "password",
			Address:  "localhost:3306",
		},
		BackendUrl:      "localhost:3000",
		OrderServiceUrl: "localhost:5001",
	}

	if dbAddress, exists := os.LookupEnv("DB_ADDRESS"); exists {
		cfg.InventoryDatabaseConfig.Address = dbAddress
	}

	if dbUser, exists := os.LookupEnv("DB_USER"); exists {
		cfg.InventoryDatabaseConfig.User = dbUser
	}

	if dbPassword, exists := os.LookupEnv("DB_PASSWORD"); exists {
		cfg.InventoryDatabaseConfig.Password = dbPassword
	}

	if dbName, exists := os.LookupEnv("DB_NAME"); exists {
		cfg.InventoryDatabaseConfig.DatabaseName = dbName
	}

	if backendUrl, exists := os.LookupEnv("BACKEND_URL"); exists {
		cfg.BackendUrl = backendUrl
	}

	if orderServiceUrl, exists := os.LookupEnv("ORDER_SERVICE_URL"); exists {
		cfg.OrderServiceUrl = orderServiceUrl
	}

	return cfg, nil
}
