package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCalculate(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   CalculateResponse
	}{
		{
			name:           "Valid simple expression",
			requestBody:    `{"expression": "2+2"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   CalculateResponse{Result: 4},
		},
		{
			name:           "Valid complex expression",
			requestBody:    `{"expression": "(2+2)*3"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   CalculateResponse{Result: 12},
		},
		{
			name:           "Invalid character",
			requestBody:    `{"expression": "2+a*2"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   CalculateResponse{Error: "expression is not valid"},
		},
		{
			name:           "Empty expression",
			requestBody:    `{"expression": ""}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   CalculateResponse{Error: "expression is not valid"},
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"wrong_key": "2+2"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   CalculateResponse{Error: "invalid request body"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Создаем ResponseRecorder для записи ответа
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			HandleCalculate(rr, req)

			// Проверяем статус ответа
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Декодируем тело ответа
			var response CalculateResponse
			err = json.NewDecoder(rr.Body).Decode(&response)
			if err != nil {
				t.Errorf("failed to decode response: %v", err)
			}

			// Проверяем тело ответа
			if tt.expectedStatus == http.StatusOK {
				if response.Result != tt.expectedBody.Result {
					t.Errorf("handler returned unexpected result: got %v want %v",
						response.Result, tt.expectedBody.Result)
				}
			} else {
				if response.Error != tt.expectedBody.Error {
					t.Errorf("handler returned unexpected error: got %v want %v",
						response.Error, tt.expectedBody.Error)
				}
			}
		})
	}

	// Тест на неподдерживаемый метод
	t.Run("Method not allowed", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/calculate", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		HandleCalculate(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code for GET request: got %v want %v",
				status, http.StatusMethodNotAllowed)
		}
	})
}
