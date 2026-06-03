import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Calculator from "./Calculator";

/**
 * Integration tests for User Story 2: input validation + error recovery.
 */
describe("Calculator — validation and error handling (US2)", () => {
  async function pressButtons(user: ReturnType<typeof userEvent.setup>, labels: string[]) {
    for (const label of labels) {
      await user.click(screen.getByRole("button", { name: label }));
    }
  }

  it("ignores a second decimal point in the same operand", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    // Press 1 . . — second . must be ignored
    await pressButtons(user, ["1", ".", "."]);
    expect(screen.getByText("1.")).toBeInTheDocument();
  });

  it("shows an error state when dividing by zero", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "÷", "0", "="]);
    // Error element must have role="alert"
    expect(screen.getByRole("alert")).toBeInTheDocument();
    expect(screen.getByRole("alert").textContent).toContain("Cannot divide by zero");
  });

  it("recovers from error state on C press", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await pressButtons(user, ["5", "÷", "0", "="]);
    expect(screen.getByRole("alert")).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "C" }));
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("ignores leading operator (operator with no prior operand)", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    // Press + before any digit — should remain at "0"
    await pressButtons(user, ["+"]);
    expect(screen.getByTestId("calculator-display")).toHaveTextContent(/^0$/);
  });

  it("last operator wins on repeated operators", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    // 5 + + * — last operator should be *
    await pressButtons(user, ["5", "+", "+", "×", "3", "="]);
    expect(screen.getByText("15")).toBeInTheDocument();
  });
});
