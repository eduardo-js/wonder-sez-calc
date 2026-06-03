import { ApiResponseError, type CalculationResult } from "./apiTypes";

const DEFAULT_BASE_URL = "http://localhost:8080";
const DEFAULT_TIMEOUT_MS = 5000;

function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL ?? DEFAULT_BASE_URL;
}

export interface CalculateOptions {
  timeoutMs?: number;
}

/**
 * POST /api/v1/calculate — delegates expression evaluation to the backend.
 * Throws ApiResponseError on non-2xx or network/timeout failures.
 */
export async function calculate(
  expression: string,
  opts: CalculateOptions = {}
): Promise<CalculationResult> {
  const { timeoutMs = DEFAULT_TIMEOUT_MS } = opts;
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);

  try {
    const response = await fetch(`${getBaseUrl()}/api/v1/calculate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ expression }),
      signal: controller.signal,
    });

    if (response.ok) {
      const data = (await response.json()) as CalculationResult;
      return data;
    }

    // Parse error envelope
    const envelope = (await response.json()) as {
      error: { code: string; message: string; fields?: Record<string, string> };
    };
    const { code, message, fields } = envelope.error;
    throw new ApiResponseError(code, message, fields);
  } catch (err) {
    if (err instanceof ApiResponseError) throw err;

    // Network or abort
    const isAbort = err instanceof Error && err.name === "AbortError";
    throw new ApiResponseError(
      "network",
      isAbort ? "Request timed out. Check your connection." : "Network error. Check your connection.",
      undefined
    );
  } finally {
    clearTimeout(timer);
  }
}
