import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
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
    expect(display).toHaveTextContent(/^2$/);
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
