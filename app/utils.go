package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alex-appy-love-story/db-lib/models/inventory"
	"github.com/alex-appy-love-story/db-lib/models/order"
	"gorm.io/gorm"
)

type OrderRequest struct {
	Username string `json:"username"`
	TokenID  uint   `json:"token_id"`
	Amount   uint   `json:"amount"`
    FailTrigger string `json:"fail_trigger"`
}


func PerformTest(db *gorm.DB, 
                 inventory []InventoryEntry, 
                 request *OrderRequest,
                 timeout int) (*order.Order, error) {
	
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Error: Could not load config")
	}

    if err := SetInventory(db, inventory); err != nil {
		return nil, fmt.Errorf("Error: Failed to set inventory: %+v", err)
	}

	if err := RequestOrder(*cfg, request); err != nil {
		return nil, fmt.Errorf("Error: %+v", err)
	}
    
	time.Sleep(time.Duration(timeout) * time.Second)

	ord, err := FetchLatestOrder(*cfg)
	if err != nil {
		return nil, fmt.Errorf("Error: %+v", err)
	}

    return ord, nil
}


func RequestOrder(cfg Config, request *OrderRequest) error {
	payloadBuf := new(bytes.Buffer)
    url := fmt.Sprintf("%s/orders", cfg.BackendUrl)
	json.NewEncoder(payloadBuf).Encode(request)
    req, _ := http.NewRequest("POST", url, payloadBuf)

	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
        return fmt.Errorf("Did not receive 200, got %+v", res.StatusCode)
	}

	return nil
}

type InventoryEntry struct {
	TokenID uint
	Amount  uint
}

func SetInventory(db *gorm.DB, info []InventoryEntry) error {
	for _, v := range info {
		if _, err := inventory.AddInventory(
			db,
			&inventory.InventoryInfo{
				TokenID: v.TokenID,
				Amount:  v.Amount,
			}); err != nil {
			return err
		}
	}

	return nil
}

func FetchLatestOrder(cfg Config) (*order.Order, error) {
    requestURL := fmt.Sprintf("http://%s", cfg.OrderServiceUrl)
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid response code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	ord := &order.Order{}

	if err := json.NewDecoder(resp.Body).Decode(&ord); err != nil {
		return nil, err
	}

	return ord, nil
}
