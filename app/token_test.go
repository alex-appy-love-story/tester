package app

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/alex-appy-love-story/db-lib/models/order"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	godotenv.Load()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.InventoryDatabaseConfig.User,
		cfg.InventoryDatabaseConfig.Password,
		cfg.InventoryDatabaseConfig.Address,
		cfg.InventoryDatabaseConfig.DatabaseName,
	)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
}

// This file contains all the tests.

func TestSuccess(t *testing.T) {
	cfg, err := LoadConfig()

	if err != nil {
		t.Errorf("Error: Could not load config")
	}

	inventory := []InventoryEntry{
		{
			TokenID: 4,
			Amount:  1,
		},
	}

	if err := SetInventory(db, inventory); err != nil {
		t.Errorf("Error: Failed to set inventory: %+v", err)
	}

	if err := RequestOrder(*cfg, &OrderRequest{
		Username: "Bob",
		TokenID:  4,
		Amount:   1,
        FailTrigger: "order",
	}); err != nil {
		t.Errorf("Error: %+v", err)
	}

	time.Sleep(1 * time.Second)

	ord, err := FetchLatestOrder(*cfg)

	if err != nil {
		t.Errorf("Error: %+v", err)
        return
	}

	if ord.OrderStatus != order.SUCCESS {
		t.Errorf("Error: Order fail")
	}
}

func TestForceFail(t *testing.T) {
    services := []string{ "order", "payment", "inventory", "delivery" }
      
    for _, service := range services {
      result, err := forceFail(t, service)
      if err != nil {
          t.Errorf("Error: " + err.Error())
          return
      }
      if result.OrderStatus != order.FORCED_FAIL {
          t.Errorf("%s failed the force fail test: returned %s", service, result.OrderStatus)
      }
    }
}

func forceFail(t *testing.T, serviceToFail string) (*order.Order, error) {
	cfg, err := LoadConfig()

	if err != nil {
		t.Errorf("Error: Could not load config")
        return &order.Order{}, fmt.Errorf("Could not load config")
	}

    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }

    if err := SetInventory(db, inventory); err != nil {
        t.Errorf("Error: Failed to set inventory: %+v", err)
        return &order.Order{}, fmt.Errorf("Failed to set inventory: %+v", err)
    }

    if err := RequestOrder(*cfg, &OrderRequest{
        Username: "Bob",
        TokenID:  4,
        Amount:   1,
        FailTrigger: serviceToFail,
    }); err != nil {
        t.Errorf("Error: %+v", err)
        return &order.Order{}, fmt.Errorf("%+v", err)
    }

    time.Sleep(5 * time.Second)

    ord, err := FetchLatestOrder(*cfg)

    if err != nil {
        t.Errorf("Error: %+v", err)
        return &order.Order{}, fmt.Errorf("%+v", err)
    }

    return ord, nil
}
