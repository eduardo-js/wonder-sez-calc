import type { Operator } from "@/hooks/calculatorTypes";

/**
 * Build an expression string from calculator operands.
 * Returns null when there is insufficient data to form an expression.
 */
export function buildExpression(
  previous: string | null,
  operator: Operator | null,
  current: string
): string | null {
  if (previous === null || operator === null) return null;
  return `${previous} ${operator} ${current}`;
}

/**
 * Build the in-progress expression text for the hint display line.
 *
 * Rules:
 * - Before any operator: show the number being typed, or "" on a fresh/overwrite state.
 * - After CHOOSE_OP (overwrite=true): show `previous operator` — RHS not typed yet.
 * - After first RHS digit (overwrite=false): show `previous operator current`.
 */
export function buildHintExpression(
  previous: string | null,
  operator: Operator | null,
  current: string,
  overwrite: boolean
): string {
  if (previous !== null && operator !== null) {
    return overwrite ? `${previous} ${operator}` : `${previous} ${operator} ${current}`;
  }
  return overwrite ? "" : current;
}
