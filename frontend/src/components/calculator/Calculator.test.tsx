import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Calculator from "./Calculator";

/**
 * Integration tests for User Story 1: basic calculation flow.
 */
describe("Calculator — basic calculation flow (US1)", () => {
  async function pressButtons(user: ReturnType<typeof userEvent.setup>, labels: string[]) {
    for (const label of labels) {
      await user.click(screen.getByRole("button", { name: label }));
    }
  }

  it("displays 19 after pressing 1 2 + 7 =", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["1", "2", "+", "7", "="]);
    expect(screen.getByText("19")).toBeInTheDocument();
  });

  it("starts a new entry (not appends) after showing a result", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "+", "3", "="]); // result: 8
    await pressButtons(user, ["2"]); // new entry
    // Should show 2, not 82
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^2$/);
  });

  it("resets display to 0 after pressing C", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["9", "+"]);
    await user.click(screen.getByRole("button", { name: "C" }));
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("registers operator and starts fresh operand", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["4", "+", "3", "="]);
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^7$/);
  });
});
