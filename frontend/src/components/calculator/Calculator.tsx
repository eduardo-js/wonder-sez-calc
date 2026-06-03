import { useCalculator } from "@/hooks/useCalculator";
import CalculatorDisplay from "./CalculatorDisplay";
import CalculatorKeypad from "./CalculatorKeypad";

/**
 * Calculator container — wires the useCalculator reducer to display + keypad.
 * Handles the full calculation, validation, and error-recovery flows.
 */
export default function Calculator() {
  const { state, handlePress } = useCalculator();

  return (
    <div className="flex w-full max-w-xs flex-col gap-3 rounded-2xl bg-slate-800 p-4 shadow-2xl sm:max-w-sm sm:gap-4 sm:p-5 md:max-w-md md:p-6">
      <CalculatorDisplay value={state.display} error={state.error} />
      <CalculatorKeypad onPress={handlePress} />
    </div>
  );
}
