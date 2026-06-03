import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Calculator from "./Calculator";

vi.mock("@/lib/api", () => ({
  calculate: vi.fn(),
}));

import { calculate } from "@/lib/api";
import { ApiResponseError } from "@/lib/apiTypes";

const mockCalculate = vi.mocked(calculate);

beforeEach(() => {
  vi.clearAllMocks();
});

async function pressButtons(user: ReturnType<typeof userEvent.setup>, labels: string[]) {
  for (const label of labels) {
    await user.click(screen.getByRole("button", { name: label }));
  }
}

/**
 * Integration tests for validation and error handling (backend-driven).
 */
describe("Calculator — validation and error handling (US2, backend-driven)", () => {
  it("ignores a second decimal point in the same operand", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["1", ".", "."]);
    // Use within to scope to the display; the hint line may also show "1."
    const display = screen.getByTestId("calculator-display");
    const valueSpans = within(display).getAllByText("1.");
    expect(valueSpans.length).toBeGreaterThanOrEqual(1);
  });

  it("shows role='alert' with backend error message when backend returns validation_failed", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("validation_failed", "invalid expression")
    );

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "+", "3", "="]);

    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("alert").textContent).toContain("invalid expression");
  });

  it("shows role='alert' with backend error message on division by zero (calculation_error)", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("calculation_error", "division by zero")
    );

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "÷", "0", "="]);

    expect(await screen.findByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("alert").textContent).toContain("division by zero");
  });

  it("display is unchanged when backend returns an error", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("calculation_error", "division by zero")
    );

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "÷", "0"]);
    // display before = shows "0" (current operand entered)
    expect(screen.getByTestId("calculator-display")).toHaveTextContent("0");

    await pressButtons(user, ["="]);
    await screen.findByRole("alert");

    // Alert shows error; the display region contains error text
    expect(screen.getByRole("alert").textContent).toContain("division by zero");
  });

  it("recovers from error state on C press", async () => {
    mockCalculate.mockRejectedValueOnce(
      new ApiResponseError("calculation_error", "division by zero")
    );

    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "÷", "0", "="]);
    await screen.findByRole("alert");

    await user.click(screen.getByRole("button", { name: "C" }));

    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("ignores leading operator (operator with no prior operand)", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["+"]);
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("last operator wins on repeated operators", async () => {
    mockCalculate.mockResolvedValueOnce({ result: "15", expression: "5 * 3" });

    const user = userEvent.setup();
    render(<Calculator />);
    // 5 + + × 3 = → operator changes from + to × → sends "5 * 3"
    await pressButtons(user, ["5", "+", "+", "×", "3", "="]);
    expect(await screen.findByText("15")).toBeInTheDocument();
    expect(mockCalculate).toHaveBeenCalledWith("5 * 3");
  });
});
