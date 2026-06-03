import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import CalculatorButton from "./CalculatorButton";

describe("CalculatorButton", () => {
  it("renders with the correct label", () => {
    render(<CalculatorButton label="7" value="7" onPress={vi.fn()} />);
    expect(screen.getByRole("button", { name: "7" })).toBeInTheDocument();
  });

  it("uses ariaLabel over label for accessible name when provided", () => {
    render(
      <CalculatorButton label="÷" value="/" onPress={vi.fn()} ariaLabel="divide" />
    );
    expect(screen.getByRole("button", { name: "divide" })).toBeInTheDocument();
  });

  it("falls back to label as accessible name when ariaLabel is not provided", () => {
    render(<CalculatorButton label="+" value="+" onPress={vi.fn()} />);
    expect(screen.getByRole("button", { name: "+" })).toBeInTheDocument();
  });

  it("calls onPress with the value exactly once on click", async () => {
    const user = userEvent.setup();
    const onPress = vi.fn();
    render(<CalculatorButton label="5" value="5" onPress={onPress} />);
    await user.click(screen.getByRole("button", { name: "5" }));
    expect(onPress).toHaveBeenCalledTimes(1);
    expect(onPress).toHaveBeenCalledWith("5");
  });

  it("applies number variant class by default", () => {
    render(<CalculatorButton label="3" value="3" onPress={vi.fn()} />);
    const btn = screen.getByRole("button", { name: "3" });
    expect(btn.className).toContain("calculator-btn-number");
  });

  it("applies operator variant class when variant='operator'", () => {
    render(<CalculatorButton label="+" value="+" variant="operator" onPress={vi.fn()} />);
    const btn = screen.getByRole("button", { name: "+" });
    expect(btn.className).toContain("calculator-btn-operator");
  });

  it("applies action variant class when variant='action'", () => {
    render(<CalculatorButton label="C" value="C" variant="action" onPress={vi.fn()} />);
    const btn = screen.getByRole("button", { name: "C" });
    expect(btn.className).toContain("calculator-btn-action");
  });
});
