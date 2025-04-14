package integrational_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	domain "pvz/internal/domain/users"
	"testing"
	"time"
)

var (
	serverURL      string
	moderatorToken string
	employeeToken  string
	pvzID          string
	receptionID    string
)

// Структуры для ответов API
type TokenResponse struct {
	Token string `json:"token"`
}

type PVZResponse struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type ReceptionResponse struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func TestMain(m *testing.M) {
	// Запуск тестового сервера (в реальном тесте нужно заменить на ваш сервер)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Здесь должен быть ваш обработчик API
		// Для примера оставлен пустым
	}))
	defer ts.Close()

	serverURL = "http://localhost:8080/api/v1"

	// Получаем токены для тестов
	if err := getTestTokens(); err != nil {
		fmt.Printf("Failed to get test tokens: %v\n", err)
		os.Exit(1)
	}

	// Запуск тестов
	code := m.Run()
	os.Exit(code)
}

func getTestTokens() error {
	// Получаем токен модератора
	moderatorData := map[string]interface{}{
		"role": domain.Moderator,
	}
	moderatorResp, err := makeRequest("POST", "/dummyLogin", moderatorData, "")
	if err != nil {
		return err
	}
	var moderatorTokenResp TokenResponse
	if err := json.Unmarshal(moderatorResp, &moderatorTokenResp); err != nil {
		return err
	}
	moderatorToken = moderatorTokenResp.Token

	// Получаем токен сотрудника
	employeeData := map[string]interface{}{
		"role": domain.Employee,
	}
	employeeResp, err := makeRequest("POST", "/dummyLogin", employeeData, "")
	if err != nil {
		return err
	}
	var employeeTokenResp TokenResponse
	if err := json.Unmarshal(employeeResp, &employeeTokenResp); err != nil {
		return err
	}

	employeeToken = employeeTokenResp.Token

	return nil
}

func makeRequest(method, path string, data interface{}, token string) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, serverURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, errorResp.Message)
	}

	return body, nil
}

func TestIntegrationFlow(t *testing.T) {
	// 1. Создаем новый ПВЗ (используем токен модератора)
	pvzData := map[string]interface{}{
		"id":               uuid.New(),
		"registrationDate": "2025-04-14T09:59:30.127Z",
		"city":             "Москва",
	}
	pvzResp, err := makeRequest("POST", "/pvz", pvzData, moderatorToken)
	if err != nil {
		t.Fatalf("Failed to create PVZ: %v", err)
	}
	var pvz PVZResponse
	if err := json.Unmarshal(pvzResp, &pvz); err != nil {
		t.Fatalf("Failed to unmarshal PVZ response: %v", err)
	}
	pvzID = pvz.ID
	t.Logf("Created PVZ with ID: %s", pvzID)

	// 2. Создаем новую приёмку (используем токен сотрудника)
	receptionData := map[string]interface{}{
		"pvzId": pvzID,
	}
	receptionResp, err := makeRequest("POST", "/receptions", receptionData, employeeToken)
	if err != nil {
		t.Fatalf("Failed to create reception: %v", err)
	}

	var reception ReceptionResponse
	if err := json.Unmarshal(receptionResp, &reception); err != nil {
		t.Fatalf("Failed to unmarshal reception response: %v", err)
	}
	receptionID = reception.ID
	t.Logf("Created reception with ID: %s", receptionID)

	// 3. Добавляем 50 товаров в приёмку
	productTypes := []string{"электроника", "одежда", "обувь"}
	for i := 0; i < 50; i++ {
		productType := productTypes[i%3]
		productData := map[string]interface{}{
			"type":  productType,
			"pvzId": pvzID,
		}
		_, err := makeRequest("POST", "/products", productData, employeeToken)
		if err != nil {
			t.Fatalf("Failed to add product %d: %v", i+1, err)
		}
	}
	t.Log("Added 50 products to the reception")

	// 4. Закрываем приёмку
	closeResp, err := makeRequest("POST", fmt.Sprintf("/pvz/%s/close_last_reception", pvzID), nil, employeeToken)
	if err != nil {
		t.Fatalf("Failed to close reception: %v", err)
	}

	var closedReception ReceptionResponse
	if err := json.Unmarshal(closeResp, &closedReception); err != nil {
		t.Fatalf("Failed to unmarshal closed reception response: %v", err)
	}

	if closedReception.Status != "close" {
		t.Errorf("Expected reception status 'close', got '%s'", closedReception.Status)
	}
	t.Logf("Reception closed successfully with status: %s", closedReception.Status)
}
