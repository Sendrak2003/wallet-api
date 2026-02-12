package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080/api/v1"

func TestIntegration_WalletOperations(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	if !isServerReady(client) {
		t.Skip("Server not available, skipping integration test")
	}

	walletID := "11111111-1111-1111-1111-111111111111"

	t.Run("GetBalance", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/wallets/" + walletID)
		if err != nil {
			t.Fatalf("failed to get balance: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if _, ok := result["balance"]; !ok {
			t.Fatal("balance field not found in response")
		}
	})

	t.Run("Deposit", func(t *testing.T) {
		payload := map[string]interface{}{
			"walletId":      walletID,
			"operationType": "DEPOSIT",
			"amount":        500,
		}

		body, _ := json.Marshal(payload)
		resp, err := client.Post(baseURL+"/wallet", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("failed to deposit: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		balance, ok := result["balance"].(float64)
		if !ok {
			t.Fatal("balance field not found or invalid type")
		}

		if balance <= 0 {
			t.Fatalf("expected positive balance, got %f", balance)
		}
	})

	t.Run("Withdraw", func(t *testing.T) {
		payload := map[string]interface{}{
			"walletId":      walletID,
			"operationType": "WITHDRAW",
			"amount":        100,
		}

		body, _ := json.Marshal(payload)
		resp, err := client.Post(baseURL+"/wallet", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("failed to withdraw: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if _, ok := result["balance"]; !ok {
			t.Fatal("balance field not found in response")
		}
	})
}

func isServerReady(client *http.Client) bool {
	resp, err := client.Get("http://localhost:8080/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
