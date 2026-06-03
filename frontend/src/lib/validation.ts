/**
 * Input validation helpers — centralized guards for calculator state transitions.
 * All functions are pure and stateless.
 */

/**
 * Returns true if `current` already contains a decimal point.
 * Used to block a second decimal input (FR-004).
 */
export function hasDecimal(current: string): boolean {
  return current.includes(".");
}

/**
 * Returns true if `current` is empty or represents only "0",
 * meaning no real operand has been entered yet.
 * Used to ignore/normalize leading operators (FR-004).
 */
export function isEmptyOperand(current: string): boolean {
  return current === "" || current === "0";
}

/**
 * Returns true if `value` is a valid arithmetic operator token.
 */
export function isOperator(value: string): boolean {
  return value === "+" || value === "-" || value === "*" || value === "/";
}
