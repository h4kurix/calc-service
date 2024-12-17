package calculator

import (
	"testing"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expected    float64
		expectErr   bool
		errorString string
	}{
		{
			name:       "Simple addition",
			expression: "2+2",
			expected:   4,
			expectErr:  false,
		},
		{
			name:       "Floating point multiplication",
			expression: "2.5*4",
			expected:   10,
			expectErr:  false,
		},
		{
			name:       "Expression with parentheses",
			expression: "(2+2)*3",
			expected:   12,
			expectErr:  false,
		},
		{
			name:        "Division by zero",
			expression:  "1/0",
			expected:    0,
			expectErr:   true,
			errorString: "деление на ноль",
		},
		{
			name:        "Invalid expression - multiple operators",
			expression:  "2++2",
			expected:    0,
			expectErr:   true,
			errorString: "недопустимое выражение",
		},
		{
			name:        "Invalid expression - operator misplacement",
			expression:  "2+*2",
			expected:    0,
			expectErr:   true,
			errorString: "недопустимое выражение",
		},
		{
			name:        "Unbalanced parentheses",
			expression:  "(2+3",
			expected:    0,
			expectErr:   true,
			errorString: "несбалансированные скобки",
		},
		{
			name:       "Complex expression",
			expression: "2*(3+4)/2",
			expected:   7,
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calc(tt.expression)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Calc(%q) expected an error, but got none", tt.expression)
					return
				}

				// Если ожидается определенный текст ошибки
				if tt.errorString != "" && err.Error() != tt.errorString {
					t.Errorf("Calc(%q) got error %v; want error with message %q",
						tt.expression, err, tt.errorString)
					return
				}
			} else {
				if err != nil {
					t.Errorf("Calc(%q) unexpected error: %v", tt.expression, err)
					return
				}

				// Сравнение с плавающей точкой с допустимой погрешностью
				if result != tt.expected {
					t.Errorf("Calc(%q) = %v; want %v", tt.expression, result, tt.expected)
				}
			}
		})
	}
}

// Тесты на токенизацию
func TestTokenize(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expected   []token
		expectErr  bool
	}{
		{
			name:       "Simple expression",
			expression: "2+3",
			expected: []token{
				{value: "2", type_: number},
				{value: "+", type_: operator},
				{value: "3", type_: number},
			},
			expectErr: false,
		},
		{
			name:       "Expression with decimal",
			expression: "2.5*4",
			expected: []token{
				{value: "2.5", type_: number},
				{value: "*", type_: operator},
				{value: "4", type_: number},
			},
			expectErr: false,
		},
		{
			name:       "Expression with parentheses",
			expression: "(2+3)*4",
			expected: []token{
				{value: "(", type_: leftParen},
				{value: "2", type_: number},
				{value: "+", type_: operator},
				{value: "3", type_: number},
				{value: ")", type_: rightParen},
				{value: "*", type_: operator},
				{value: "4", type_: number},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenize(tt.expression)

			if tt.expectErr {
				if err == nil {
					t.Errorf("tokenize(%q) expected an error, but got none", tt.expression)
				}
				return
			}

			if err != nil {
				t.Errorf("tokenize(%q) unexpected error: %v", tt.expression, err)
				return
			}

			// Проверка токенов
			if len(tokens) != len(tt.expected) {
				t.Errorf("tokenize(%q) returned %d tokens; want %d",
					tt.expression, len(tokens), len(tt.expected))
				return
			}

			for i, token := range tokens {
				if token.value != tt.expected[i].value || token.type_ != tt.expected[i].type_ {
					t.Errorf("tokenize(%q) token %d: got %v; want %v",
						tt.expression, i, token, tt.expected[i])
				}
			}
		})
	}
}
