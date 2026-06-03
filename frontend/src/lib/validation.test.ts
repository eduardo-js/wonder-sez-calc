import { describe, it, expect } from "vitest";
import { hasDecimal, isEmptyOperand, isOperator } from "./validation";

describe("hasDecimal", () => {
  it("returns false when no decimal", () => {
    expect(hasDecimal("123")).toBe(false);
  });

  it("returns true when decimal present", () => {
    expect(hasDecimal("12.3")).toBe(true);
  });

  it("returns true for a bare decimal", () => {
    expect(hasDecimal(".")).toBe(true);
  });
});

describe("isEmptyOperand", () => {
  it("returns true for empty string", () => {
    expect(isEmptyOperand("")).toBe(true);
  });

  it("returns true for '0'", () => {
    expect(isEmptyOperand("0")).toBe(true);
  });

  it("returns false for any non-zero digit", () => {
    expect(isEmptyOperand("5")).toBe(false);
  });

  it("returns false for a number starting with non-zero", () => {
    expect(isEmptyOperand("12")).toBe(false);
  });
});

describe("isOperator", () => {
  it("returns true for +", () => expect(isOperator("+")).toBe(true));
  it("returns true for -", () => expect(isOperator("-")).toBe(true));
  it("returns true for *", () => expect(isOperator("*")).toBe(true));
  it("returns true for /", () => expect(isOperator("/")).toBe(true));
  it("returns false for digit", () => expect(isOperator("5")).toBe(false));
  it("returns false for =", () => expect(isOperator("=")).toBe(false));
  it("returns false for C", () => expect(isOperator("C")).toBe(false));
});
