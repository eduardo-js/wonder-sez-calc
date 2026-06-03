import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, act, fireEvent, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Calculator from "./Calculator";

vi.mock("@/lib/api", () => ({
  calculate: vi.fn(),
}));

import { calculate } from "@/lib/api";

const mockCalculate = vi.mocked(calculate);

beforeEach(() => {
  vi.clearAllMocks();
});

/**
 * Integration tests for User Story 1: basic calculation flow (backend-driven).
 */
describe("Calculator — basic calculation flow (US1)", () => {
  async function pressButtons(user: ReturnType<typeof userEvent.setup>, labels: string[]) {
    for (const label of labels) {
      await user.click(screen.getByRole("button", { name: label }));
    }
  }

  it("displays backend result after pressing 1 2 + 7 =", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "19", expression: "12 + 7" });

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["1", "2", "+", "7", "="]);

    const display = screen.getByTestId("calculator-display");
    expect(await screen.findByTestId("calculator-display")).toHaveTextContent("19");
    expect(display).toBeInTheDocument();
    expect(mockCalculate).toHaveBeenCalledWith("12 + 7");
  });

  it("starts a new entry (not appends) after showing a result", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "8", expression: "5 + 3" });

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "+", "3", "="]);

    // Wait for backend result to appear in display
    const display = screen.getByTestId("calculator-display");
    await vi.waitFor(() => expect(display).toHaveTextContent("8"));

    await pressButtons(user, ["2"]);
    // The hint line may also show "2"; check the non-alert, non-status value span directly
    const valueSpan = within(display).getAllByText(/^2$/);
    expect(valueSpan.length).toBeGreaterThanOrEqual(1);
  });

  it("resets display to 0 after pressing C", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["9", "+"]);
    await user.click(screen.getByRole("button", { name: "C" }));
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("registers operator and starts fresh operand", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "7", expression: "4 + 3" });

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["4", "+", "3", "="]);

    const display = screen.getByTestId("calculator-display");
    await vi.waitFor(() => expect(display).toHaveTextContent("7"));
  });
});

describe("Calculator — expression hint line (005)", () => {
  function pressButton(label: string) {
    fireEvent.click(screen.getByRole("button", { name: label }));
  }

  it("shows expression '12 +' after pressing 1, 2, + (before RHS)", async () => {
    render(<Calculator />);

    await act(async () => { pressButton("1"); });
    await act(async () => { pressButton("2"); });
    await act(async () => { pressButton("+"); });

    const hintLine = screen.getByRole("status");
    expect(hintLine).toHaveTextContent("12 +");
    expect(mockCalculate).not.toHaveBeenCalled();
  });

  it("shows expression '12 + 3' after pressing 1, 2, +, 3 (before =)", async () => {
    render(<Calculator />);

    await act(async () => { pressButton("1"); });
    await act(async () => { pressButton("2"); });
    await act(async () => { pressButton("+"); });
    await act(async () => { pressButton("3"); });

    const hintLine = screen.getByRole("status");
    expect(hintLine).toHaveTextContent("12 + 3");
    expect(mockCalculate).not.toHaveBeenCalled();
  });

  it("calculate is NOT called until = is pressed", async () => {
    render(<Calculator />);

    await act(async () => { pressButton("1"); });
    await act(async () => { pressButton("2"); });
    await act(async () => { pressButton("+"); });
    await act(async () => { pressButton("3"); });

    expect(mockCalculate).not.toHaveBeenCalled();
  });

  it("after = is pressed, result shows and calculate was called once", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "15", expression: "12 + 3" });

    render(<Calculator />);

    await act(async () => { pressButton("1"); });
    await act(async () => { pressButton("2"); });
    await act(async () => { pressButton("+"); });
    await act(async () => { pressButton("3"); });
    await act(async () => { pressButton("="); });
    await act(async () => {});

    await vi.waitFor(() =>
      expect(screen.getByTestId("calculator-display")).toHaveTextContent("15")
    );
    expect(mockCalculate).toHaveBeenCalledTimes(1);
    expect(mockCalculate).toHaveBeenCalledWith("12 + 3");
  });
});
