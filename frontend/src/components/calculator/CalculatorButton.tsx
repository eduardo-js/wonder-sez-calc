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
  /** When true, renders a spinner in place of label and disables the button */
  isLoading?: boolean;
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
 * When isLoading=true (only used for '='), shows a spinner and disables interaction.
 */
export default function CalculatorButton({
  label,
  value,
  variant = "number",
  onPress,
  ariaLabel,
  className,
  isLoading = false,
}: CalculatorButtonProps) {
  return (
    <Button
      variant="ghost"
      aria-label={isLoading ? "calculating" : (ariaLabel ?? label)}
      disabled={isLoading}
      className={cn(
        "h-14 w-full rounded-lg text-xl font-semibold sm:h-16 md:h-20",
        "focus-visible:ring-2 focus-visible:ring-sky-400 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-800 focus-visible:outline-none",
        variantClassMap[variant],
        className
      )}
      onClick={() => !isLoading && onPress(value)}
    >
      {isLoading ? (
        <span
          data-testid="spinner"
          aria-hidden="true"
          className="inline-block h-5 w-5 animate-spin rounded-full border-2 border-white border-t-transparent"
        />
      ) : (
        label
      )}
    </Button>
  );
}
