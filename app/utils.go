package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alex-appy-love-story/db-lib/models/inventory"
	"github.com/alex-appy-love-story/db-lib/models/order"
	"gorm.io/gorm"
)

type OrderRequest struct {
	Username string `json:"username"`
	TokenID  uint   `json:"token_id"`
	Amount   uint   `json:"amount"`
    FailTrigger string
}

func RequestOrder(cfg Config, request *OrderRequest) error {
	payloadBuf := new(bytes.Buffer)
    url := fmt.Sprintf("http://%s/orders", cfg.BackendUrl)
	json.NewEncoder(payloadBuf).Encode(request)
    req, _ := http.NewRequest("POST", url, payloadBuf)

	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Did not receive 200")
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
