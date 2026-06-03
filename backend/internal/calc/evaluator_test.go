package calc_test

import (
	"errors"
	"testing"

	"github.com/wonderlic-calc/backend/internal/calc"
)

func TestEvaluate_HappyPath(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		// Basic arithmetic
		{"addition", "2 + 3", "5"},
		{"subtraction", "10 - 4", "6"},
		{"multiplication", "3 * 4", "12"},
		{"division", "10 / 4", "2.5"},

		// Operator precedence
		{"precedence mul over add", "2 + 3 * 4", "14"},
		{"precedence mul over sub", "10 - 2 * 3", "4"},
		{"precedence div over add", "8 / 4 + 1", "3"},

		// Parentheses override precedence
		{"parens override", "(2 + 3) * 4", "20"},
		{"nested parens", "((2 + 3)) * 4", "20"},
		{"parens sub", "10 / (2 + 3)", "2"},

		// Decimals
		{"decimal operands", "1.5 + 2.5", "4"},
		{"decimal result", "1.0 / 3.0", "0.333333333333"},
		{"float formatting", "0.1 + 0.2", "0.3"},

		// Unary minus
		{"unary minus at start", "-3 + 5", "2"},
		{"unary minus after paren", "(-3 + 5)", "2"},
		{"unary minus after operator", "5 + -3", "2"},
		{"unary minus multiply", "-3 * -2", "6"},

		// Integer results
		{"integer plain", "6 / 2", "3"},
		{"negative result", "3 - 10", "-7"},
		{"zero", "0 + 0", "0"},

		// Chained operations
		{"left assoc sub", "10 - 3 - 2", "5"},
		{"left assoc div", "12 / 4 / 3", "1"},
		{"complex", "3 + 4 * 2 / (1 - 5)", "1"},

		// Formatting: ≤12 significant digits, trailing zeros trimmed
		{"trim trailing zeros", "1.10 + 0", "1.1"},
		{"large integer", "1000000 * 1000000", "1000000000000"},

		// Large-magnitude results: scientific notation, never int64-corrupted (US1/US2)
		{"reported overflow regression", "-9223372036854775808 - 999999999999999999999999999999999999999999", "-1e+42"},
		{"exact-int boundary stays plain", "9007199254740992", "9007199254740992"},
		{"just past boundary goes scientific", "9007199254740992 + 9007199254740992", "1.80143985095e+16"},
		{"large positive scientific", "1e40 * 2", "2e+40"},
		{"very small non-zero not zero", "1 / 1e20", "1e-20"},

		// Scientific notation in number literals
		{"sci notation", "1e3 + 0", "1000"},
		{"sci notation with plus sign", "1e+3 + 0", "1000"},
		{"sci notation with minus sign", "1e-3 + 0", "0.001"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := calc.Evaluate(tc.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("Evaluate(%q) = %q, want %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestEvaluate_Errors(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantErr error
	}{
		// Invalid expression
		{"empty string", "", calc.ErrInvalidExpression},
		{"only spaces", "   ", calc.ErrInvalidExpression},
		{"consecutive operators", "5 + * 2", calc.ErrInvalidExpression},
		{"missing right operand", "5 +", calc.ErrInvalidExpression},
		{"missing left operand", "+ 5", calc.ErrInvalidExpression},
		{"unbalanced open paren", "(3 + 4", calc.ErrInvalidExpression},
		{"unbalanced close paren", "3 + 4)", calc.ErrInvalidExpression},
		{"only operator", "*", calc.ErrInvalidExpression},
		{"only unary minus", "-", calc.ErrInvalidExpression},
		{"two numbers no op", "2 3", calc.ErrInvalidExpression},
		{"double decimal", "1.2.3 + 1", calc.ErrInvalidExpression},
		{"letters", "abc", calc.ErrInvalidExpression},
		{"empty parens", "()", calc.ErrInvalidExpression},

		// Calculation errors
		{"divide by zero", "1 / 0", calc.ErrDivideByZero},
		{"divide by zero expr", "(3 - 3) * 0 / (1 - 1)", calc.ErrDivideByZero},

		// Non-finite
		{"overflow large", "1e308 * 1e308", calc.ErrNonFinite},
		{"over-large literal", "1e400", calc.ErrNonFinite},

		// Scientific notation edge cases
		{"invalid sci notation no digits", "1e + 2", calc.ErrInvalidExpression},
		{"standalone e", "e + 1", calc.ErrInvalidExpression},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := calc.Evaluate(tc.expr)
			if err == nil {
				t.Fatalf("expected error %v, got nil", tc.wantErr)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}
