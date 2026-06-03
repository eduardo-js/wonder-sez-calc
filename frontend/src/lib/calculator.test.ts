import { describe, it, expect } from "vitest";
import { compute } from "./calculator";

describe("compute", () => {
  it("adds two numbers", () => {
    expect(compute("12", "+", "7")).toEqual({ value: "19", error: null });
  });

  it("subtracts two numbers", () => {
    expect(compute("10", "-", "3")).toEqual({ value: "7", error: null });
  });

  it("multiplies two numbers", () => {
    expect(compute("3", "*", "4")).toEqual({ value: "12", error: null });
  });

  it("divides two numbers", () => {
    expect(compute("10", "/", "2")).toEqual({ value: "5", error: null });
  });

  it("returns decimal result for non-integer division", () => {
    const result = compute("1", "/", "3");
    expect(result.error).toBeNull();
    expect(parseFloat(result.value)).toBeCloseTo(0.333333, 5);
  });

  it("returns error on divide by zero", () => {
    const result = compute("5", "/", "0");
    expect(result.error).toBe("Cannot divide by zero");
    expect(result.value).toBe("0");
  });

  it("returns error on invalid operand", () => {
    const result = compute("abc", "+", "1");
    expect(result.error).toBe("Invalid operand");
  });

  it("handles negative results", () => {
    expect(compute("3", "-", "10")).toEqual({ value: "-7", error: null });
  });

  it("handles float addition", () => {
    const result = compute("0.1", "+", "0.2");
    expect(result.error).toBeNull();
    // Floating point formatted to 12 significant digits
    expect(parseFloat(result.value)).toBeCloseTo(0.3, 10);
  });
});
