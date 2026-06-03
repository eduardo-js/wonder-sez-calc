import { useReducer } from "react";
import { compute } from "../lib/calculator";
import { hasDecimal, isEmptyOperand } from "../lib/validation";
import {
  type CalculatorAction,
  type CalculatorState,
  INITIAL_STATE,
} from "./calculatorTypes";

/**
 * Pure reducer implementing all state transitions per data-model.md.
 */
export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction
): CalculatorState {
  // Any action on error state is ignored except CLEAR
  if (state.error && action.type !== "CLEAR") {
    return state;
  }

  switch (action.type) {
    case "INPUT_DIGIT": {
      const { digit } = action;

      if (state.overwrite) {
        // Start fresh — replace display entirely
        const newCurrent = digit === "0" ? "0" : digit;
        return { ...state, current: newCurrent, display: newCurrent, overwrite: false };
      }

      // Normalize leading zeros: don't allow "007"
      if (state.current === "0" && digit === "0") return state;
      if (state.current === "0" && digit !== ".") {
        const newCurrent = digit;
        return { ...state, current: newCurrent, display: newCurrent };
      }

      const newCurrent = state.current + digit;
      return { ...state, current: newCurrent, display: newCurrent };
    }

    case "INPUT_DECIMAL": {
      // Ignore if already has decimal
      if (hasDecimal(state.current)) return state;

      if (state.overwrite) {
        // Start fresh decimal: "0."
        return { ...state, current: "0.", display: "0.", overwrite: false };
      }

      const newCurrent = state.current + ".";
      return { ...state, current: newCurrent, display: newCurrent };
    }

    case "CHOOSE_OP": {
      const { operator } = action;

      // Leading operator: no previous and current is "0" with overwrite → ignore
      if (state.previous === null && isEmptyOperand(state.current) && state.overwrite) {
        return state;
      }

      // Repeated/changed operator: just swap the operator, keep previous
      if (state.previous !== null && state.overwrite) {
        return { ...state, operator };
      }

      // Normal case: commit current → previous, set operator, set overwrite for next operand
      return {
        ...state,
        previous: state.current,
        operator,
        overwrite: true,
        display: state.current,
      };
    }

    case "EVALUATE": {
      // Nothing to evaluate without all parts
      if (state.previous === null || state.operator === null) {
        return state;
      }

      const result = compute(state.previous, state.operator, state.current);

      if (result.error) {
        return {
          ...state,
          error: result.error,
          display: result.error,
        };
      }

      return {
        ...state,
        display: result.value,
        current: result.value,
        previous: null,
        operator: null,
        overwrite: true,
        error: null,
      };
    }

    case "CLEAR": {
      return { ...INITIAL_STATE };
    }

    default: {
      // TypeScript exhaustiveness check — unreachable at runtime
      throw new Error(`Unhandled action: ${JSON.stringify(action)}`);
    }
  }
}

/**
 * useCalculator — wires the reducer and returns state + dispatch helpers.
 */
export function useCalculator() {
  const [state, dispatch] = useReducer(calculatorReducer, INITIAL_STATE);

  function handlePress(value: string): void {
    if (value === "C") {
      dispatch({ type: "CLEAR" });
      return;
    }
    if (value === "=") {
      dispatch({ type: "EVALUATE" });
      return;
    }
    if (value === ".") {
      dispatch({ type: "INPUT_DECIMAL" });
      return;
    }
    if (["+", "-", "*", "/"].includes(value)) {
      dispatch({ type: "CHOOSE_OP", operator: value as "+" | "-" | "*" | "/" });
      return;
    }
    // Digit
    dispatch({ type: "INPUT_DIGIT", digit: value });
  }

  return { state, handlePress };
}
