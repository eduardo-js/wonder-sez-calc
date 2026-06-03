/**
 * Single source of truth for the calculator keypad layout.
 * Each entry drives one CalculatorButton via CalculatorKeypad.
 */

export type ButtonKind = "digit" | "operator" | "decimal" | "equals" | "clear";
export type ButtonVariant = "number" | "operator" | "action";

export interface ButtonConfig {
  /** Stable key for rendering/tests */
  id: string;
  /** Visible text on the button */
  label: string;
  /** Logical token emitted on press */
  value: string;
  /** Drives reducer action + visual variant */
  kind: ButtonKind;
  /** Maps to shadcn Button styling */
  variant: ButtonVariant;
  /** Grid column span (default 1) */
  span?: number;
}

export const KEYPAD_BUTTONS: ButtonConfig[] = [
  { id: "clear", label: "C", value: "C", kind: "clear", variant: "action", span: 2 },
  { id: "divide", label: "÷", value: "/", kind: "operator", variant: "operator" },
  { id: "multiply", label: "×", value: "*", kind: "operator", variant: "operator" },

  { id: "digit-7", label: "7", value: "7", kind: "digit", variant: "number" },
  { id: "digit-8", label: "8", value: "8", kind: "digit", variant: "number" },
  { id: "digit-9", label: "9", value: "9", kind: "digit", variant: "number" },
  { id: "subtract", label: "−", value: "-", kind: "operator", variant: "operator" },

  { id: "digit-4", label: "4", value: "4", kind: "digit", variant: "number" },
  { id: "digit-5", label: "5", value: "5", kind: "digit", variant: "number" },
  { id: "digit-6", label: "6", value: "6", kind: "digit", variant: "number" },
  { id: "add", label: "+", value: "+", kind: "operator", variant: "operator" },

  { id: "digit-1", label: "1", value: "1", kind: "digit", variant: "number" },
  { id: "digit-2", label: "2", value: "2", kind: "digit", variant: "number" },
  { id: "digit-3", label: "3", value: "3", kind: "digit", variant: "number" },
  { id: "equals", label: "=", value: "=", kind: "equals", variant: "action" },

  { id: "digit-0", label: "0", value: "0", kind: "digit", variant: "number", span: 2 },
  { id: "decimal", label: ".", value: ".", kind: "decimal", variant: "number" },
];
