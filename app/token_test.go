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
    fmt.Println("Testing Success")
	inventory := []InventoryEntry{
		{
			TokenID: 4,
			Amount:  1,
		},
	}

    request := &OrderRequest{
		Username: "Bob",
		TokenID:  4,
		Amount:   1,
    }

    ord, err := PerformTest(db, inventory, request, 4)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

	if ord.OrderStatus != order.SUCCESS {
        t.Errorf("Error: Order is not success, status is %s", ord.OrderStatus)
        return
	}
}


//
// Testing all fail types: Out of stock, insufficienct funds, token not found
//
func TestOutOfStock(t *testing.T) {
    fmt.Println("Testing out of stock")
    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }

    request := &OrderRequest{
        Username: "Boba",
        TokenID:  1,
        Amount:   2,
    }

    ord, err := PerformTest(db, inventory, request, 3)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if ord.OrderStatus != order.INVENTORY_FAIL_STOCK {
		t.Errorf("Error: Order did not fail correctly, status is %s", ord.OrderStatus)
        return
	}
}

func TestInsufficientFunds(t *testing.T) {
    fmt.Println("Testing insufficient funds")
    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  30,
        },
    }

    request := &OrderRequest{
        Username: "Bobby",
        TokenID:  4,
        Amount:   90,
    }

    ord, err := PerformTest(db, inventory, request, 3)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if ord.OrderStatus != order.PAYMENT_FAIL_INSUFFICIENT {
		t.Errorf("Error: Order did not fail correctly, status is %s", ord.OrderStatus)
        return
	}
}


func TestTokenNotFound(t *testing.T) {
    fmt.Println("Testing token not found")
    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }

    request := &OrderRequest{
        Username: "Bobi",
        TokenID:  42,
        Amount:   1,
    }

    ord, err := PerformTest(db, inventory, request, 3)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if ord.OrderStatus != order.PAYMENT_FAIL_TOKEN_NOT_FOUND {
		t.Errorf("Error: Order did not fail correctly, status is %s", ord.OrderStatus)
        return
	}
}


//
// Testing force fails
//
func TestForceFails(t *testing.T) {
    fmt.Println("Testing force fails")
    
    fmt.Println("--- Sending success request")

    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }
    successRequest := &OrderRequest{
        Username: "Bobo",
        TokenID:  4,
        Amount:   1,
    }
    
    _, err := PerformTest(db, inventory, successRequest, 4)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    
    services := []string{ "order", "payment", "inventory", "delivery" }

    for _, service := range services {
        fmt.Println("--- Testing", service)
        result, err := forceFail(t, service)
        if err != nil {
            t.Errorf("Error: " + err.Error())
            return
        }
        if result.OrderStatus != order.FORCED_FAIL {
            t.Errorf("%s failed the force fail test: returned %s", service, result.OrderStatus)
            return
        }
    }
}

func forceFail(t *testing.T, serviceToFail string) (*order.Order, error) {
    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }

    request := &OrderRequest{
        Username: "Bobo",
        TokenID:  4,
        Amount:   1,
        FailTrigger: serviceToFail,
    }

    ord, err := PerformTest(db, inventory, request, 3)
    if err != nil {
        return &order.Order{}, err
    }

    return ord, nil
}


//
// Testing circuit breaker
//
func TestCircuitBreaker(t *testing.T) {
    fmt.Println("Testing circuit breaker")
    
    inventory := []InventoryEntry{
        {
            TokenID: 4,
            Amount:  1,
        },
    }

    successRequest := &OrderRequest{
        Username: "Bobo",
        TokenID:  4,
        Amount:   1,
    }
    failRequest := &OrderRequest{
        Username: "Bobo",
        TokenID:  4,
        Amount:   1,
        FailTrigger: "order",
    }

    fmt.Println("--- Sending success request")
    ord, err := PerformTest(db, inventory, successRequest, 4)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    for i := 0; i < 9; i++ {
        fmt.Println("--- Sending fail request")
        ord, err = PerformTest(db, inventory, failRequest, 3)
        if err != nil {
            t.Errorf(err.Error())
            return
        }
    }

    time.Sleep(4 * time.Second)

    if ord.OrderStatus != order.DEFAULT_RESPONSE {
        t.Errorf("Failed the circuit breaker test: returned %s", ord.OrderStatus)
        return
    }

}
