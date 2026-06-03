import { cn } from "@/lib/utils";

export interface CalculatorDisplayProps {
  /** Current display string */
  value: string;
  /** When set, renders an accessible error alert */
  error: string | null;
  /** Backend/network error message (shown in role="alert") */
  errorMsg?: string | null;
  /** In-progress expression text, e.g. "12 +" or "12 + 3"; omit or "" to hide */
  expression?: string;
}

/**
 * Calculator display area. Shows value normally; renders error with role="alert"
 * when error or errorMsg is set. Renders a secondary expression hint line above
 * the main value when `expression` is a non-empty string and no error is shown.
 */
export default function CalculatorDisplay({
  value,
  error,
  errorMsg,
  expression,
}: CalculatorDisplayProps) {
  const alertText = errorMsg ?? error ?? null;
  const showExpression = !alertText && expression != null && expression !== "";

  return (
    <div
      data-testid="calculator-display"
      className={cn(
        "flex min-h-[4rem] w-full flex-col items-end justify-end overflow-hidden rounded-lg px-4 py-3 sm:min-h-[5rem] md:min-h-[6rem]",
        alertText ? "bg-red-900 text-red-200" : "bg-slate-900 text-white"
      )}
    >
      {showExpression && (
        <span
          role="status"
          aria-live="polite"
          className="mb-1 max-w-full truncate text-right font-mono text-sm text-slate-400"
        >
          {expression}
        </span>
      )}
      {alertText ? (
        <span
          role="alert"
          aria-live="assertive"
          className="max-w-full truncate text-right text-lg font-medium sm:text-xl md:text-2xl"
        >
          {alertText}
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
