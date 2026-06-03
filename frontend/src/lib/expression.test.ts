import { describe, it, expect } from "vitest";
import { buildExpression, buildHintExpression } from "./expression";

describe("buildExpression", () => {
  it("returns null when previous is null", () => {
    expect(buildExpression(null, "+", "3")).toBeNull();
  });

  it("returns null when operator is null", () => {
    expect(buildExpression("12", null, "3")).toBeNull();
  });

  it("builds a full expression string", () => {
    expect(buildExpression("12", "+", "3")).toBe("12 + 3");
  });
});

describe("buildHintExpression", () => {
  it("typing first number, not overwrite → shows the number", () => {
    expect(buildHintExpression(null, null, "12", false)).toBe("12");
  });

  it("fresh state, overwrite=true → empty string (suppress redundant '0')", () => {
    expect(buildHintExpression(null, null, "0", true)).toBe("");
  });

  it("after operator pressed, overwrite=true → shows 'previous operator' only", () => {
    expect(buildHintExpression("12", "+", "12", true)).toBe("12 +");
  });

  it("after RHS typed, overwrite=false → shows 'previous operator current'", () => {
    expect(buildHintExpression("12", "+", "3", false)).toBe("12 + 3");
  });

  it("works with subtraction operator", () => {
    expect(buildHintExpression("5", "-", "2", false)).toBe("5 - 2");
  });

  it("works with multiplication operator", () => {
    expect(buildHintExpression("7", "*", "8", false)).toBe("7 * 8");
  });

  it("works with division operator, overwrite=true (no RHS)", () => {
    expect(buildHintExpression("9", "/", "9", true)).toBe("9 /");
  });
});
