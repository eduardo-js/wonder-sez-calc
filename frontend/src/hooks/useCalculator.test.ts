import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useCalculator } from "./useCalculator";

// Mock the api module
vi.mock("@/lib/api", () => ({
  calculate: vi.fn(),
}));

import { calculate } from "@/lib/api";
import { ApiResponseError } from "@/lib/apiTypes";

const mockCalculate = vi.mocked(calculate);

beforeEach(() => {
  vi.resetAllMocks();
  vi.useFakeTimers();
});

afterEach(() => {
  vi.useRealTimers();
});

describe("useCalculator — EVALUATE (backend-driven)", () => {
  it("calls calculate() and sets display from backend result", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "14", expression: "2 + 3 * 4" });

    const { result } = renderHook(() => useCalculator());

    await act(async () => {
      result.current.handlePress("2");
    });
    await act(async () => {
      result.current.handlePress("+");
    });
    await act(async () => {
      result.current.handlePress("3");
    });
    await act(async () => {
      await result.current.handlePress("=");
    });

    expect(mockCalculate).toHaveBeenCalledOnce();
    expect(result.current.state.display).toBe("14");
    expect(result.current.state.status).toBe("idle");
    expect(result.current.state.errorMsg).toBeNull();
  });

  it("second EVALUATE while loading issues no request", async () => {
    let resolveCalc!: (v: { result: string; expression: string }) => void;
    mockCalculate.mockReturnValueOnce(
      new Promise((res) => {
        resolveCalc = res;
      })
    );

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("5"); });
    await act(async () => { result.current.handlePress("+"); });
    await act(async () => { result.current.handlePress("3"); });

    // First press — starts loading (don't await full resolution)
    act(() => {
      void result.current.handlePress("=");
    });

    // Wait for loading state
    await act(async () => {});

    expect(result.current.state.status).toBe("loading");

    // Second press while loading — must NOT issue another request
    await act(async () => {
      await result.current.handlePress("=");
    });

    expect(mockCalculate).toHaveBeenCalledTimes(1);

    // Resolve the first request
    await act(async () => {
      resolveCalc({ result: "8", expression: "5 + 3" });
    });

    expect(result.current.state.status).toBe("idle");
  });

  it("status toggles loading → idle on successful resolve", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "9", expression: "4 + 5" });

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("4"); });
    await act(async () => { result.current.handlePress("+"); });
    await act(async () => { result.current.handlePress("5"); });
    await act(async () => { await result.current.handlePress("="); });

    expect(result.current.state.status).toBe("idle");
    expect(result.current.state.display).toBe("9");
  });

  it("sets status error and errorMsg on backend error; display unchanged", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("calculation_error", "division by zero")
    );

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("5"); });
    await act(async () => { result.current.handlePress("/"); });
    await act(async () => { result.current.handlePress("0"); });

    const displayBefore = result.current.state.display;

    await act(async () => { await result.current.handlePress("="); });

    expect(result.current.state.status).toBe("error");
    expect(result.current.state.errorMsg).toBe("division by zero");
    expect(result.current.state.display).toBe(displayBefore);
  });

  it("clears errorMsg and status on next digit input", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("validation_failed", "invalid expression")
    );

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("5"); });
    await act(async () => { result.current.handlePress("+"); });
    await act(async () => { await result.current.handlePress("="); });

    expect(result.current.state.status).toBe("error");

    await act(async () => { result.current.handlePress("3"); });

    expect(result.current.state.status).toBe("idle");
    expect(result.current.state.errorMsg).toBeNull();
  });

  it("clears errorMsg on CLEAR", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("network", "Network error. Check your connection.")
    );

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("1"); });
    await act(async () => { result.current.handlePress("+"); });
    await act(async () => { result.current.handlePress("1"); });
    await act(async () => { await result.current.handlePress("="); });

    expect(result.current.state.status).toBe("error");

    await act(async () => { result.current.handlePress("C"); });

    expect(result.current.state.status).toBe("idle");
    expect(result.current.state.errorMsg).toBeNull();
  });

  it("typing digits and operators does NOT call calculate; only = does", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "5", expression: "2 + 3" });

    const { result } = renderHook(() => useCalculator());

    await act(async () => { result.current.handlePress("2"); });
    await act(async () => { result.current.handlePress("+"); });
    await act(async () => { result.current.handlePress("3"); });

    // No call yet
    expect(mockCalculate).not.toHaveBeenCalled();

    await act(async () => { await result.current.handlePress("="); });

    expect(mockCalculate).toHaveBeenCalledTimes(1);
  });
});
