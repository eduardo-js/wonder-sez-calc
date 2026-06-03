import { useReducer } from "react";
import { calculate } from "../lib/api";
import { ApiResponseError } from "../lib/apiTypes";
import { buildExpression } from "../lib/expression";
import { hasDecimal, isEmptyOperand } from "../lib/validation";
import {
  type CalculatorAction,
  type CalculatorState,
  INITIAL_STATE,
} from "./calculatorTypes";

/**
 * Pure reducer implementing all state transitions.
 * EVALUATE is now handled async in the hook; reducer handles lifecycle actions.
 */
export function calculatorReducer(
  state: CalculatorState,
  action: CalculatorAction
): CalculatorState {
  switch (action.type) {
    case "INPUT_DIGIT": {
      const { digit } = action;

      // Clear error state on new input
      const clearError =
        state.status === "error"
          ? { status: "idle" as const, errorMsg: null, error: null }
          : {};

      if (state.overwrite) {
        const newCurrent = digit === "0" ? "0" : digit;
        return {
          ...state,
          ...clearError,
          current: newCurrent,
          display: newCurrent,
          overwrite: false,
        };
      }

      if (state.current === "0" && digit === "0") return { ...state, ...clearError };
      if (state.current === "0" && digit !== ".") {
        const newCurrent = digit;
        return { ...state, ...clearError, current: newCurrent, display: newCurrent };
      }

      const newCurrent = state.current + digit;
      return { ...state, ...clearError, current: newCurrent, display: newCurrent };
    }

    case "INPUT_DECIMAL": {
      const clearError =
        state.status === "error"
          ? { status: "idle" as const, errorMsg: null, error: null }
          : {};

      if (hasDecimal(state.current)) return { ...state, ...clearError };

      if (state.overwrite) {
        return {
          ...state,
          ...clearError,
          current: "0.",
          display: "0.",
          overwrite: false,
        };
      }

      const newCurrent = state.current + ".";
      return { ...state, ...clearError, current: newCurrent, display: newCurrent };
    }

    case "CHOOSE_OP": {
      const { operator } = action;

      const clearError =
        state.status === "error"
          ? { status: "idle" as const, errorMsg: null, error: null }
          : {};

      if (state.previous === null && isEmptyOperand(state.current) && state.overwrite) {
        return state;
      }

      if (state.previous !== null && state.overwrite) {
        return { ...state, ...clearError, operator };
      }

      return {
        ...state,
        ...clearError,
        previous: state.current,
        operator,
        overwrite: true,
        display: state.current,
      };
    }

    case "EVALUATE": {
      // Sync EVALUATE is now a no-op: async logic handled in hook.
      // Guard: do nothing if no expression parts.
      if (state.previous === null || state.operator === null) return state;
      // Guard: do nothing if already loading
      if (state.status === "loading") return state;
      return state;
    }

    case "EVALUATE_START": {
      return { ...state, status: "loading", errorMsg: null };
    }

    case "EVALUATE_SUCCESS": {
      return {
        ...state,
        status: "idle",
        errorMsg: null,
        error: null,
        display: action.result,
        current: action.result,
        previous: null,
        operator: null,
        overwrite: true,
      };
    }

    case "EVALUATE_ERROR": {
      return {
        ...state,
        status: "error",
        errorMsg: action.errorMsg,
        // display is intentionally unchanged
      };
    }

    case "CLEAR": {
      return { ...INITIAL_STATE };
    }

    default: {
      throw new Error(`Unhandled action: ${JSON.stringify(action)}`);
    }
  }
}

/**
 * useCalculator — wires the reducer and returns state + dispatch helpers.
 * EVALUATE is async: builds expression, calls backend, dispatches lifecycle actions.
 */
export function useCalculator() {
  const [state, dispatch] = useReducer(calculatorReducer, INITIAL_STATE);

  async function handlePress(value: string): Promise<void> {
    if (value === "C") {
      dispatch({ type: "CLEAR" });
      return;
    }
    if (value === "=") {
      // Guard: single in-flight
      if (state.status === "loading") return;
      // Guard: need both operands
      if (state.previous === null || state.operator === null) return;

      const expression = buildExpression(state.previous, state.operator, state.current);
      if (expression === null) return;

      dispatch({ type: "EVALUATE_START" });
      try {
        const result = await calculate(expression);
        dispatch({ type: "EVALUATE_SUCCESS", result: result.result });
      } catch (err) {
        const msg =
          err instanceof ApiResponseError
            ? err.message
            : "An unexpected error occurred.";
        dispatch({ type: "EVALUATE_ERROR", errorMsg: msg });
      }
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
