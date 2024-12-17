package calculator

import (
	"strconv"
	"strings"
	"unicode"

	"calc-service/internal/errors"
	"calc-service/pkg/logger"
)

// Calc вычисляет результат математического выражения
func Calc(expression string) (float64, error) {
	// Токенизация выражения
	tokens, err := tokenize(expression)
	if err != nil {
		logger.Error("Tokenization error: %v", err)
		return 0, err
	}

	// Вычисление выражения
	result, err := evaluateExpression(tokens)
	if err != nil {
		logger.Error("Expression evaluation failed: %v", err)
		return 0, err
	}

	logger.Info("Expression successfully evaluated: %s = %f", expression, result)
	return result, nil
}

// Структура для хранения токена
type token struct {
	value string
	type_ tokenType
}

type tokenType int

const (
	number tokenType = iota
	operator
	leftParen
	rightParen
)

// Токенизация выражения
func tokenize(expression string) ([]token, error) {
	var tokens []token
	var currentNumber strings.Builder
	var parenCount int // Счётчик скобок

	for i := 0; i < len(expression); i++ {
		ch := expression[i]

		switch {
		case unicode.IsDigit(rune(ch)) || ch == '.':
			currentNumber.WriteByte(ch)
			if i == len(expression)-1 || !(unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				// Проверка числа на формат (недопустимые точки)
				if strings.Count(currentNumber.String(), ".") > 1 {
					return nil, errors.ErrInvalidExpression
				}
				tokens = append(tokens, token{currentNumber.String(), number})
				currentNumber.Reset()
			}

		case ch == '+' || ch == '-' || ch == '*' || ch == '/':
			if currentNumber.Len() > 0 {
				tokens = append(tokens, token{currentNumber.String(), number})
				currentNumber.Reset()
			}
			tokens = append(tokens, token{string(ch), operator})

		case ch == '(':
			tokens = append(tokens, token{"(", leftParen})
			parenCount++ // Увеличиваем счётчик открывающих скобок

		case ch == ')':
			if currentNumber.Len() > 0 {
				tokens = append(tokens, token{currentNumber.String(), number})
				currentNumber.Reset()
			}
			tokens = append(tokens, token{")", rightParen})
			parenCount-- // Уменьшаем счётчик скобок
			if parenCount < 0 {
				return nil, errors.ErrUnbalancedParentheses // Если скобки несбалансированы
			}

		default:
			return nil, errors.ErrInvalidSymbol
		}
	}

	if parenCount != 0 {
		return nil, errors.ErrUnbalancedParentheses // Проверяем на конце, если есть незакрытые скобки
	}

	return tokens, nil
}

// Применение операции
func applyOp(a float64, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, errors.ErrUnknownOperator
	}
}

// Вспомогательная функция для обработки операций
func popAndApply(values *[]float64, ops *[]string) error {
	if len(*values) < 2 {
		return errors.ErrInvalidExpression
	}
	b, a := (*values)[len(*values)-1], (*values)[len(*values)-2]
	*values = (*values)[:len(*values)-2]

	op := (*ops)[len(*ops)-1]
	*ops = (*ops)[:len(*ops)-1]

	result, err := applyOp(a, b, op)
	if err != nil {
		return err
	}
	*values = append(*values, result)
	return nil
}

// Вычисление выражения
func evaluateExpression(tokens []token) (float64, error) {
	var values []float64
	var ops []string

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		switch t.type_ {
		case number:
			val, err := strconv.ParseFloat(t.value, 64)
			if err != nil {
				return 0, err
			}
			values = append(values, val)

		case leftParen:
			ops = append(ops, t.value)

		case rightParen:
			for len(ops) > 0 && ops[len(ops)-1] != "(" {
				if err := popAndApply(&values, &ops); err != nil {
					return 0, err
				}
			}
			if len(ops) == 0 {
				return 0, errors.ErrUnbalancedParentheses
			}
			ops = ops[:len(ops)-1] // Удаляем открывающую скобку

		case operator:
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(t.value) {
				if err := popAndApply(&values, &ops); err != nil {
					return 0, err
				}
			}
			ops = append(ops, t.value)
		}
	}

	for len(ops) > 0 {
		if err := popAndApply(&values, &ops); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, errors.ErrInvalidExpression
	}

	return values[0], nil
}

// Приоритет операторов
func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
