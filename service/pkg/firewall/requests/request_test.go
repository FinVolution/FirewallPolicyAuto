package requests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock HTTP Server
func setupMockServer() *httptest.Server {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"message": "created"}`))
		} else if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message": "success"}`))
		}
	}))
	return mockServer
}

func TestRequest(t *testing.T) {
	mockServer := setupMockServer()
	defer mockServer.Close()

	hc := NewHTTPClient(false)

	// Test GET method
	getParams := RequestParams{
		URL:    mockServer.URL,
		Method: http.MethodGet,
	}
	response, err := hc.Request(&getParams)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", response.StatusCode)
	}
	var getResponse map[string]string
	if err := json.Unmarshal(response.Body, &getResponse); err != nil {
		t.Fatalf("Failed to unmarshal GET response body: %v", err)
	}
	if getResponse["message"] != "success" {
		t.Fatalf("Expected message 'success', got '%s'", getResponse["message"])
	}

	// Test POST method
	postParams := RequestParams{
		URL:    mockServer.URL,
		Method: http.MethodPost,
		Body:   struct{ Message string }{Message: "test"},
	}
	response, err = hc.Request(&postParams)
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status Created, got %d", response.StatusCode)
	}
	var postResponse map[string]string
	if err := json.Unmarshal(response.Body, &postResponse); err != nil {
		t.Fatalf("Failed to unmarshal POST response body: %v", err)
	}
	if postResponse["message"] != "created" {
		t.Fatalf("Expected message 'created', got '%s'", postResponse["message"])
	}
}
