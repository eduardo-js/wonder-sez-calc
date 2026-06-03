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

describe("CalculatorDisplay — expression hint line (005)", () => {
  it("renders expression as a secondary polite live region", () => {
    render(
      <CalculatorDisplay
        value="3"
        error={null}
        expression="12 +"
      />
    );

    const hintLine = screen.getByRole("status");
    expect(hintLine).toHaveTextContent("12 +");
    expect(hintLine).toHaveAttribute("aria-live", "polite");
  });

  it("renders expression with RHS typed", () => {
    render(
      <CalculatorDisplay
        value="3"
        error={null}
        expression="12 + 3"
      />
    );

    const hintLine = screen.getByRole("status");
    expect(hintLine).toHaveTextContent("12 + 3");
  });

  it("empty expression string suppresses hint line", () => {
    render(
      <CalculatorDisplay
        value="0"
        error={null}
        expression=""
      />
    );

    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("omitted expression suppresses hint line", () => {
    render(
      <CalculatorDisplay
        value="5"
        error={null}
      />
    );

    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("suppresses hint line when error is set", () => {
    render(
      <CalculatorDisplay
        value="3"
        error="invalid input"
        expression="12 + 3"
      />
    );

    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("suppresses hint line when errorMsg is set", () => {
    render(
      <CalculatorDisplay
        value="3"
        error={null}
        errorMsg="division by zero"
        expression="12 + 3"
      />
    );

    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });
});
