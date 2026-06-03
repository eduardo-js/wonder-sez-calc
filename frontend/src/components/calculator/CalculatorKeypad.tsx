import { KEYPAD_BUTTONS } from "@/lib/keypad";
import CalculatorButton from "./CalculatorButton";

export interface CalculatorKeypadProps {
  /** Forwarded to each CalculatorButton */
  onPress: (value: string) => void;
}

/**
 * Renders one CalculatorButton per KEYPAD_BUTTONS entry in a responsive CSS Grid.
 * 4-column layout; buttons with span=2 occupy two columns.
 */
export default function CalculatorKeypad({ onPress }: CalculatorKeypadProps) {
  return (
    <div
      className="grid grid-cols-4 gap-2 sm:gap-3 md:gap-4"
      data-testid="calculator-keypad"
    >
      {KEYPAD_BUTTONS.map((btn) => (
        <CalculatorButton
          key={btn.id}
          label={btn.label}
          value={btn.value}
          variant={btn.variant}
          onPress={onPress}
          className={btn.span === 2 ? "col-span-2" : undefined}
        />
      ))}
    </div>
  );
}
