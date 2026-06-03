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
