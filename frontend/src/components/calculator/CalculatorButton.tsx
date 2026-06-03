import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export type ButtonVariant = "number" | "operator" | "action";

export interface CalculatorButtonProps {
  /** Rendered text content */
  label: string;
  /** Token passed to onPress */
  value: string;
  /** Visual variant — maps to distinct styling */
  variant?: ButtonVariant;
  /** Called once per activation with the button value */
  onPress: (value: string) => void;
  /** Accessible name; falls back to label */
  ariaLabel?: string;
  /** Optional extra classes (e.g. col-span-2) */
  className?: string;
}

const variantClassMap: Record<ButtonVariant, string> = {
  number:
    "bg-slate-700 text-white hover:bg-slate-600 active:bg-slate-500 calculator-btn-number",
  operator:
    "bg-amber-500 text-white hover:bg-amber-400 active:bg-amber-300 calculator-btn-operator",
  action:
    "bg-slate-400 text-white hover:bg-slate-300 active:bg-slate-200 calculator-btn-action",
};

/**
 * Reusable calculator button wrapping the shadcn Button primitive.
 * Applies variant-specific styling and forwards the value on press.
 */
export default function CalculatorButton({
  label,
  value,
  variant = "number",
  onPress,
  ariaLabel,
  className,
}: CalculatorButtonProps) {
  return (
    <Button
      variant="ghost"
      aria-label={ariaLabel ?? label}
      className={cn(
        "h-14 w-full rounded-lg text-xl font-semibold sm:h-16 md:h-20",
        variantClassMap[variant],
        className
      )}
      onClick={() => onPress(value)}
    >
      {label}
    </Button>
  );
}
