package handler

import (
	"calc-service/internal/calculator"
	"calc-service/internal/errors"
	"encoding/json"
	"net/http"
	"strings"
)

// CalculateRequest represents the structure of the incoming JSON request
type CalculateRequest struct {
	Expression string `json:"expression"`
}

// CalculateResponse represents the structure of the response
type CalculateResponse struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// ValidateExpression checks if the expression contains only valid characters
func validateExpression(expression string) error {
	validChars := "0123456789.+-*/()"
	for _, ch := range expression {
		if !strings.ContainsRune(validChars, ch) {
			return errors.ErrInvalidSymbol
		}
	}
	return nil
}

// HandleCalculate handles the calculation request
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		sendErrorResponse(w, errors.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	// Decode the request
	var req CalculateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Отключает обработку незнакомых полей
	if err := decoder.Decode(&req); err != nil {
		sendErrorResponse(w, errors.ErrInvalidRequestBody.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Удаляем пробелы из выражения
	req.Expression = strings.ReplaceAll(req.Expression, " ", "")

	// Проверка на пустое выражение или недопустимые символы
	if req.Expression == "" || validateExpression(req.Expression) != nil {
		sendErrorResponse(w, errors.ErrInvalidSymbol.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Perform calculation
	result, err := calculator.Calc(req.Expression)
	if err != nil {
		// For internal errors, send 500 Internal Server Error
		sendErrorResponse(w, errors.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CalculateResponse{Result: result})
}

// sendErrorResponse is a helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(CalculateResponse{Error: errorMsg})
}
