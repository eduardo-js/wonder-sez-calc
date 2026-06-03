import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import CalculatorDisplay from "./CalculatorDisplay";

describe("CalculatorDisplay", () => {
  it("renders a scientific-notation result fully and legibly (US2)", () => {
    render(<CalculatorDisplay value="-1e+42" error={null} />);
    const display = screen.getByTestId("calculator-display");
    // The full scientific-notation string is present, not clipped or reformatted.
    expect(display).toHaveTextContent("-1e+42");
    // No alert styling for a valid value.
    expect(screen.queryByRole("alert")).not.toBeInTheDocument();
  });

  it("renders a very small scientific-notation result", () => {
    render(<CalculatorDisplay value="1e-20" error={null} />);
    expect(screen.getByTestId("calculator-display")).toHaveTextContent("1e-20");
  });

  it("surfaces an evaluation error in an alert region without showing the value (US3)", () => {
    render(
      <CalculatorDisplay
        value="123"
        error={null}
        errorMsg="result is not a finite number"
      />
    );
    const alert = screen.getByRole("alert");
    expect(alert).toHaveTextContent("result is not a finite number");
    // Prior value is replaced by the alert, not corrupted alongside it.
    expect(alert).not.toHaveTextContent("123");
  });
});
