// Package calc provides a pure mathematical expression evaluator.
// It supports +, -, *, /, parentheses, decimal numbers, and unary minus.
// No HTTP or external dependencies.
package calc

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Sentinel errors returned by Evaluate.
var (
	// ErrInvalidExpression is returned when the expression cannot be parsed
	// (empty, malformed, unbalanced parentheses, illegal characters, etc.).
	ErrInvalidExpression = errors.New("invalid expression")

	// ErrDivideByZero is returned when the expression attempts division by zero.
	ErrDivideByZero = errors.New("division by zero")

	// ErrNonFinite is returned when evaluation produces Inf or NaN.
	ErrNonFinite = errors.New("result is not a finite number")
)

// tokenKind classifies a lexer token.
type tokenKind int

const (
	kindNumber   tokenKind = iota
	kindOperator           // + - * / u-
	kindLParen             // (
	kindRParen             // )
)

// token is a single unit produced by the lexer.
type token struct {
	kind  tokenKind
	value string // operator symbol or number literal
}

// Evaluate tokenizes, parses, and evaluates the mathematical expression expr.
// It returns the formatted result string or a sentinel error.
func Evaluate(expr string) (string, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return "", err
	}
	rpn, err := toRPN(tokens)
	if err != nil {
		return "", err
	}
	result, err := evalRPN(rpn)
	if err != nil {
		return "", err
	}
	return formatResult(result), nil
}

// tokenize converts expr into a slice of tokens.
// Unary minus is detected and represented as the token "u-".
func tokenize(expr string) ([]token, error) {
	trimmed := strings.TrimSpace(expr)
	if trimmed == "" {
		return nil, fmt.Errorf("%w: expression is empty", ErrInvalidExpression)
	}

	// Validate allowed characters: digits, '.', '+', '-', '*', '/', '(', ')', spaces,
	// and 'e'/'E' (part of float literals like 1e308).
	for _, ch := range trimmed {
		if !unicode.IsDigit(ch) && ch != '.' && ch != '+' && ch != '-' &&
			ch != '*' && ch != '/' && ch != '(' && ch != ')' && ch != ' ' &&
			ch != 'e' && ch != 'E' {
			return nil, fmt.Errorf("%w: illegal character %q", ErrInvalidExpression, ch)
		}
	}

	var tokens []token
	runes := []rune(trimmed)
	n := len(runes)
	i := 0

	for i < n {
		ch := runes[i]

		// Skip spaces
		if ch == ' ' {
			i++
			continue
		}

		// Number: digit, or leading dot, or scientific notation handled as part of number
		if unicode.IsDigit(ch) || ch == '.' {
			j := i
			dotCount := 0
			for j < n && (unicode.IsDigit(runes[j]) || runes[j] == '.') {
				if runes[j] == '.' {
					dotCount++
					if dotCount > 1 {
						return nil, fmt.Errorf("%w: invalid number with multiple decimal points", ErrInvalidExpression)
					}
				}
				j++
			}
			// Optional scientific notation: e/E followed by optional +/- and digits
			if j < n && (runes[j] == 'e' || runes[j] == 'E') {
				j++ // consume e/E
				if j < n && (runes[j] == '+' || runes[j] == '-') {
					j++ // consume sign
				}
				start := j
				for j < n && unicode.IsDigit(runes[j]) {
					j++
				}
				if j == start {
					return nil, fmt.Errorf("%w: invalid scientific notation", ErrInvalidExpression)
				}
			}
			num := string(runes[i:j])
			tokens = append(tokens, token{kind: kindNumber, value: num})
			i = j
			continue
		}

		// Parentheses
		if ch == '(' {
			tokens = append(tokens, token{kind: kindLParen, value: "("})
			i++
			continue
		}
		if ch == ')' {
			tokens = append(tokens, token{kind: kindRParen, value: ")"})
			i++
			continue
		}

		// Operators: detect unary minus
		if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
			isUnary := false
			if ch == '-' {
				if len(tokens) == 0 {
					isUnary = true
				} else {
					prev := tokens[len(tokens)-1]
					if prev.kind == kindOperator || prev.kind == kindLParen {
						isUnary = true
					}
				}
			}
			if isUnary {
				tokens = append(tokens, token{kind: kindOperator, value: "u-"})
			} else {
				tokens = append(tokens, token{kind: kindOperator, value: string(ch)})
			}
			i++
			continue
		}

		// 'e'/'E' appearing outside a number context (e.g. standalone "e")
		return nil, fmt.Errorf("%w: unexpected character %q", ErrInvalidExpression, ch)
	}

	return tokens, nil
}

// precedence returns operator precedence (higher = tighter binding).
// Callers only ever pass +, -, *, /, u- so the switch is exhaustive.
func precedence(op string) int {
	switch op {
	case "*", "/":
		return 2
	case "u-":
		return 3
	default: // +, -
		return 1
	}
}

// toRPN converts tokens to Reverse Polish Notation using the shunting-yard algorithm.
// Returns ErrInvalidExpression for structural problems (unbalanced parens, etc.).
func toRPN(tokens []token) ([]token, error) {
	var output []token
	var opStack []token

	push := func(t token) { opStack = append(opStack, t) }
	pop := func() token {
		top := opStack[len(opStack)-1]
		opStack = opStack[:len(opStack)-1]
		return top
	}
	peek := func() token { return opStack[len(opStack)-1] }

	for _, tok := range tokens {
		switch tok.kind {
		case kindNumber:
			output = append(output, tok)

		case kindOperator:
			// Unary minus is right-associative: only pop operators with strictly higher precedence.
			// Binary operators are left-associative: pop operators with equal or higher precedence.
			isRightAssoc := tok.value == "u-"
			for len(opStack) > 0 && peek().kind == kindOperator {
				topPrec := precedence(peek().value)
				curPrec := precedence(tok.value)
				if (!isRightAssoc && topPrec >= curPrec) || (isRightAssoc && topPrec > curPrec) {
					output = append(output, pop())
				} else {
					break
				}
			}
			push(tok)

		case kindLParen:
			push(tok)

		case kindRParen:
			// Pop until matching '('
			foundLParen := false
			for len(opStack) > 0 {
				top := peek()
				if top.kind == kindLParen {
					pop() // discard the '('
					foundLParen = true
					break
				}
				output = append(output, pop())
			}
			if !foundLParen {
				return nil, fmt.Errorf("%w: unmatched closing parenthesis", ErrInvalidExpression)
			}
		}
	}

	// Flush remaining operators; any leftover LParen means unmatched open paren
	for len(opStack) > 0 {
		top := pop()
		if top.kind == kindLParen {
			return nil, fmt.Errorf("%w: unmatched opening parenthesis", ErrInvalidExpression)
		}
		output = append(output, top)
	}

	return output, nil
}

// evalRPN evaluates an RPN token sequence and returns the numeric result.
func evalRPN(rpn []token) (float64, error) {
	if len(rpn) == 0 {
		return 0, fmt.Errorf("%w: empty expression", ErrInvalidExpression)
	}

	var stack []float64

	push := func(v float64) { stack = append(stack, v) }
	pop := func() (float64, bool) {
		if len(stack) == 0 {
			return 0, false
		}
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v, true
	}

	for _, tok := range rpn {
		switch tok.kind {
		case kindNumber:
			v, _ := strconv.ParseFloat(tok.value, 64) // validated during tokenization
			push(v)

		case kindOperator:
			if tok.value == "u-" {
				a, ok := pop()
				if !ok {
					return 0, fmt.Errorf("%w: missing operand for unary minus", ErrInvalidExpression)
				}
				push(-a)
				continue
			}

			b, okB := pop()
			a, okA := pop()
			if !okA || !okB {
				return 0, fmt.Errorf("%w: not enough operands for %q", ErrInvalidExpression, tok.value)
			}

			var result float64
			switch tok.value {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			default: // "/"
				if b == 0 {
					return 0, fmt.Errorf("%w", ErrDivideByZero)
				}
				result = a / b
			}
			push(result)
		}
	}

	if len(stack) != 1 {
		return 0, fmt.Errorf("%w: malformed expression", ErrInvalidExpression)
	}

	result := stack[0]
	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, fmt.Errorf("%w", ErrNonFinite)
	}

	return result, nil
}

// maxExactInt is the largest magnitude at which every consecutive integer is
// exactly representable in a float64 (2^53 = 9007199254740992). Beyond it, plain
// integer rendering is misleading and an int64 cast would overflow/saturate, so
// such results are formatted in scientific notation instead.
const maxExactInt = 1 << 53

// formatResult formats a float64 result per the contract:
//   - whole numbers within the exact-integer range (|v| ≤ 2^53) rendered plain
//     (no decimal point)
//   - larger magnitudes and non-integers (incl. very small values) to ≤12
//     significant digits in scientific notation, trailing zeros stripped
func formatResult(v float64) string {
	if v == math.Trunc(v) && !math.IsInf(v, 0) && math.Abs(v) <= maxExactInt {
		return strconv.FormatInt(int64(v), 10)
	}
	// %g with 12 significant digits, trailing zeros stripped automatically
	return strconv.FormatFloat(v, 'g', 12, 64)
}
