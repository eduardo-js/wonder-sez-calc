import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import CalculatorKeypad from "./CalculatorKeypad";
import CalculatorButton from "./CalculatorButton";

/**
 * Responsive structure tests (US3).
 * jsdom does not apply CSS breakpoints, so tests verify the presence
 * of responsive Tailwind classes in the DOM class strings.
 */
describe("Calculator responsive structure (US3)", () => {
  it("keypad grid has responsive gap classes", () => {
    render(<CalculatorKeypad onPress={vi.fn()} />);
    const grid = screen.getByTestId("calculator-keypad");
    // gap-2 at base, sm:gap-3, md:gap-4
    expect(grid.className).toContain("gap-2");
    expect(grid.className).toContain("sm:gap-3");
    expect(grid.className).toContain("md:gap-4");
  });

  it("CalculatorButton has responsive height classes", () => {
    render(<CalculatorButton label="7" value="7" onPress={vi.fn()} />);
    const btn = screen.getByRole("button", { name: "7" });
    // h-14 base, sm:h-16 minimum
    expect(btn.className).toContain("h-14");
    expect(btn.className).toContain("sm:h-16");
  });

  it("zero button has col-span-2 class", () => {
    render(<CalculatorKeypad onPress={vi.fn()} />);
    const zeroBtn = screen.getByRole("button", { name: "0" });
    expect(zeroBtn.className).toContain("col-span-2");
  });

  it("clear button has col-span-2 class", () => {
    render(<CalculatorKeypad onPress={vi.fn()} />);
    const clearBtn = screen.getByRole("button", { name: "C" });
    expect(clearBtn.className).toContain("col-span-2");
  });
});
