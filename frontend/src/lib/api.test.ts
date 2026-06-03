import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { calculate } from "./api";
import { ApiResponseError } from "./apiTypes";

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

function makeJsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

describe("calculate()", () => {
  it("returns CalculationResult on 200", async () => {
    mockFetch.mockResolvedValueOnce(
      makeJsonResponse({ result: "14", expression: "2 + 3 * 4" })
    );
    const result = await calculate("2 + 3 * 4");
    expect(result.result).toBe("14");
    expect(result.expression).toBe("2 + 3 * 4");
  });

  it("sends POST to /api/v1/calculate with JSON body", async () => {
    mockFetch.mockResolvedValueOnce(
      makeJsonResponse({ result: "5", expression: "2 + 3" })
    );
    await calculate("2 + 3");
    expect(mockFetch).toHaveBeenCalledOnce();
    const [url, init] = mockFetch.mock.calls[0] as [string, RequestInit];
    expect(url).toContain("/api/v1/calculate");
    expect(init.method).toBe("POST");
    expect(JSON.parse(init.body as string)).toEqual({ expression: "2 + 3" });
  });

  it("throws ApiResponseError with code and message on 400 validation_failed", async () => {
    mockFetch.mockResolvedValueOnce(
      makeJsonResponse(
        {
          error: {
            code: "validation_failed",
            message: "invalid expression",
            fields: { expression: "unexpected token" },
          },
        },
        400
      )
    );
    let caught: ApiResponseError | null = null;
    try {
      await calculate("5 + * 2");
    } catch (err) {
      caught = err as ApiResponseError;
    }
    expect(caught).toBeInstanceOf(ApiResponseError);
    expect(caught?.code).toBe("validation_failed");
    expect(caught?.message).toBe("invalid expression");
    expect(caught?.fields?.expression).toBe("unexpected token");
  });

  it("throws ApiResponseError with code calculation_error on 422", async () => {
    mockFetch.mockResolvedValueOnce(
      makeJsonResponse(
        { error: { code: "calculation_error", message: "division by zero" } },
        422
      )
    );
    let caught: ApiResponseError | null = null;
    try {
      await calculate("1 / 0");
    } catch (err) {
      caught = err as ApiResponseError;
    }
    expect(caught).toBeInstanceOf(ApiResponseError);
    expect(caught?.code).toBe("calculation_error");
    expect(caught?.message).toBe("division by zero");
  });

  it("maps AbortError to connectivity ApiResponseError with code 'network'", async () => {
    const abortErr = new Error("aborted");
    abortErr.name = "AbortError";
    mockFetch.mockRejectedValueOnce(abortErr);

    let caught: ApiResponseError | null = null;
    try {
      await calculate("1 + 1", { timeoutMs: 50000 });
    } catch (err) {
      caught = err as ApiResponseError;
    }
    expect(caught).toBeInstanceOf(ApiResponseError);
    expect(caught?.code).toBe("network");
  });

  it("maps network failure to connectivity error", async () => {
    mockFetch.mockRejectedValueOnce(new TypeError("Failed to fetch"));
    let caught: ApiResponseError | null = null;
    try {
      await calculate("1 + 1");
    } catch (err) {
      caught = err as ApiResponseError;
    }
    expect(caught).toBeInstanceOf(ApiResponseError);
    expect(caught?.code).toBe("network");
    expect(caught?.message).toMatch(/network|connection/i);
  });
});
