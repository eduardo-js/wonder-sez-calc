import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import CalculatorKeypad from "./CalculatorKeypad";
import { KEYPAD_BUTTONS } from "@/lib/keypad";

describe("CalculatorKeypad", () => {
  it("renders exactly one button per KEYPAD_BUTTONS entry", () => {
    render(<CalculatorKeypad onPress={vi.fn()} />);
    const buttons = screen.getAllByRole("button");
    expect(buttons).toHaveLength(KEYPAD_BUTTONS.length);
  });

  it("renders a grid container", () => {
    render(<CalculatorKeypad onPress={vi.fn()} />);
    const grid = screen.getByTestId("calculator-keypad");
    expect(grid.className).toContain("grid");
    expect(grid.className).toContain("grid-cols-4");
  });

  it("passes correct value to onPress for each button", async () => {
    const onPress = vi.fn();
    const { getByRole } = render(<CalculatorKeypad onPress={onPress} />);
    // Test a few buttons by accessible name
    const clearBtn = getByRole("button", { name: "C" });
    clearBtn.click();
    expect(onPress).toHaveBeenCalledWith("C");
  });
});
