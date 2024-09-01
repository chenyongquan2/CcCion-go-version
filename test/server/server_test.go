package server

import (
	"CcCoin-go-version/internal/blockchain"
	"CcCoin-go-version/internal/encryption"
	"CcCoin-go-version/internal/server"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBlockchainServer_AddTransaction(t *testing.T) {
	mockBlockchain := blockchain.NewBlockchain(3)
	server := server.NewBlockchainServer(mockBlockchain)
	senderPrivateKey, senderPublicKey := encryption.GenerateKeyPair()
	_, receiverPublicKey := encryption.GenerateKeyPair()

	testCases := []struct {
		name           string
		txData         map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid Transaction",
			txData: map[string]interface{}{
				"SenderPublicKey":   senderPublicKey,
				"SenderPrivateKey":  senderPrivateKey,
				"ReceiverPublicKey": receiverPublicKey,
				"Amount":            100.0,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid Transaction Data",
			txData: map[string]interface{}{
				"InvalidField": "invalidValue",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.txData)
			req, _ := http.NewRequest("POST", "/transction/", bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}

func TestBlockchainServer_StartMineTask(t *testing.T) {
	mockBlockchain := blockchain.NewBlockchain(3)
	server := server.NewBlockchainServer(mockBlockchain)

	testCases := []struct {
		name           string
		mineData       map[string]string
		expectedStatus int
	}{
		{
			name: "Valid Mine Request",
			mineData: map[string]string{
				"MinerPublicKey": "minerPublicKey",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid Mine Data",
			mineData:       map[string]string{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.mineData)
			req, _ := http.NewRequest("POST", "/mine/", bytes.NewBuffer(jsonData))
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}

func TestBlockchainServer_TransactionHandler(t *testing.T) {
	mockBlockchain := blockchain.NewBlockchain(3)
	server := server.NewBlockchainServer(mockBlockchain)

	testCases := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET Method",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         "PUT",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/transction/", nil)
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}

func TestBlockchainServer_MineHandler(t *testing.T) {
	mockBlockchain := blockchain.NewBlockchain(3)
	server := server.NewBlockchainServer(mockBlockchain)

	testCases := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "GET Method",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method",
			method:         "PUT",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.method, "/mine/", nil)
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}
		})
	}
}
