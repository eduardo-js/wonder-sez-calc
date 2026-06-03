import { cn } from "@/lib/utils";

export interface CalculatorDisplayProps {
  /** Current display string */
  value: string;
  /** When set, renders an error state with role="alert" */
  error: string | null;
}

/**
 * Calculator display area. Shows value normally; renders error with role="alert"
 * and distinct styling when error is set (FR-005 / contract).
 */
export default function CalculatorDisplay({ value, error }: CalculatorDisplayProps) {
  return (
    <div
      data-testid="calculator-display"
      className={cn(
        "flex min-h-[4rem] w-full items-end justify-end overflow-hidden rounded-lg px-4 py-3 sm:min-h-[5rem] md:min-h-[6rem]",
        error ? "bg-red-900 text-red-200" : "bg-slate-900 text-white"
      )}
    >
      {error ? (
        <span
          role="alert"
          className="max-w-full truncate text-right text-lg font-medium sm:text-xl md:text-2xl"
        >
          {error}
        </span>
      ) : (
        <span
          className={cn(
            "max-w-full truncate text-right font-mono font-semibold",
            value.length > 12
              ? "text-lg sm:text-xl"
              : value.length > 8
                ? "text-2xl sm:text-3xl"
                : "text-3xl sm:text-4xl md:text-5xl"
          )}
        >
          {value}
        </span>
      )}
    </div>
  );
}
