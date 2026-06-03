/**
 * CalculatorState shape and action union — per data-model.md.
 */

export type Operator = "+" | "-" | "*" | "/";

export interface CalculatorState {
  /** What the user sees */
  display: string;
  /** Operand currently being typed */
  current: string;
  /** Stored left-hand operand */
  previous: string | null;
  /** Pending operation */
  operator: Operator | null;
  /** If true, next digit replaces display (after = or operator) */
  overwrite: boolean;
  /** Set on invalid operation (e.g. divide-by-zero); blocks compute */
  error: string | null;
}

export type CalculatorAction =
  | { type: "INPUT_DIGIT"; digit: string }
  | { type: "INPUT_DECIMAL" }
  | { type: "CHOOSE_OP"; operator: Operator }
  | { type: "EVALUATE" }
  | { type: "CLEAR" };

export const INITIAL_STATE: CalculatorState = {
  display: "0",
  current: "0",
  previous: null,
  operator: null,
  overwrite: true,
  error: null,
};
