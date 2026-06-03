/**
 * Pure arithmetic computation. No side effects.
 */

export type Operator = "+" | "-" | "*" | "/";

export interface ComputeResult {
  value: string;
  error: string | null;
}

/**
 * Compute `left op right` and return a formatted result or an error.
 * Divide-by-zero returns `{ value: "0", error: "Cannot divide by zero" }`.
 */
export function compute(left: string, operator: Operator, right: string): ComputeResult {
  const a = parseFloat(left);
  const b = parseFloat(right);

  if (isNaN(a) || isNaN(b)) {
    return { value: "0", error: "Invalid operand" };
  }

  if (operator === "/" && b === 0) {
    return { value: "0", error: "Cannot divide by zero" };
  }

  let result: number;
  switch (operator) {
    case "+":
      result = a + b;
      break;
    case "-":
      result = a - b;
      break;
    case "*":
      result = a * b;
      break;
    case "/":
      result = a / b;
      break;
    default: {
      const _: never = operator;
      throw new Error(`Unhandled operator: ${String(_)}`);
    }
  }

  // Format: avoid floating-point noise, cap at 12 significant digits
  const formatted =
    Number.isInteger(result) ? result.toString() : parseFloat(result.toPrecision(12)).toString();

  return { value: formatted, error: null };
}
