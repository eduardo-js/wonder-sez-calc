import { test, expect, type Page } from "@playwright/test";

// Button accessible names match the visible labels in keypad.ts.
// Operators render as Unicode glyphs on the keypad (÷ × − +).
const OP = { divide: "÷", multiply: "×", subtract: "−", add: "+" } as const;

async function press(page: Page, ...labels: string[]) {
  for (const label of labels) {
    await page.getByRole("button", { name: label, exact: true }).click();
  }
}

function display(page: Page) {
  return page.getByTestId("calculator-display");
}

test.beforeEach(async ({ page }) => {
  await page.goto("/");
});

test("loads with an initial display of 0", async ({ page }) => {
  await expect(display(page)).toContainText("0");
});

test("addition: 2 + 2 = 4", async ({ page }) => {
  await press(page, "2", OP.add, "2", "=");
  await expect(display(page)).toContainText("4");
});

test("subtraction: 9 − 4 = 5", async ({ page }) => {
  await press(page, "9", OP.subtract, "4", "=");
  await expect(display(page)).toContainText("5");
});

test("multiplication: 6 × 7 = 42", async ({ page }) => {
  await press(page, "6", OP.multiply, "7", "=");
  await expect(display(page)).toContainText("42");
});

test("division: 10 ÷ 4 = 2.5", async ({ page }) => {
  await press(page, "1", "0", OP.divide, "4", "=");
  await expect(display(page)).toContainText("2.5");
});

test("decimal operands: 1.5 + 2.25 = 3.75", async ({ page }) => {
  await press(page, "1", ".", "5", OP.add, "2", ".", "2", "5", "=");
  await expect(display(page)).toContainText("3.75");
});

test("chained operations reuse the previous result", async ({ page }) => {
  await press(page, "2", OP.add, "3", "="); // 5
  await expect(display(page)).toContainText("5");
  await press(page, OP.multiply, "4", "="); // 5 * 4 = 20
  await expect(display(page)).toContainText("20");
});

test("clear resets the display to 0", async ({ page }) => {
  await press(page, "5", "5");
  await expect(display(page)).toContainText("55");
  await press(page, "C");
  await expect(display(page)).toContainText("0");
});

test("shows the in-progress expression hint", async ({ page }) => {
  await press(page, "1", "2", OP.add);
  // Hint renders the logical operator token (+), not the keypad glyph.
  await expect(page.getByRole("status")).toHaveText("12 +");
  await press(page, "3");
  await expect(page.getByRole("status")).toHaveText("12 + 3");
});

test("division by zero surfaces a backend error alert", async ({ page }) => {
  await press(page, "5", OP.divide, "0", "=");
  const alert = page.getByRole("alert");
  await expect(alert).toBeVisible();
  await expect(alert).toContainText(/division by zero/i);
});

test("recovers from an error on the next input", async ({ page }) => {
  await press(page, "5", OP.divide, "0", "=");
  await expect(page.getByRole("alert")).toBeVisible();
  await press(page, "7");
  await expect(page.getByRole("alert")).toHaveCount(0);
  await expect(display(page)).toContainText("7");
});

test("large-number computation round-trips through the backend", async ({ page }) => {
  // 999999999 * 999999999 -> backend formats large results in scientific notation.
  await press(page, ...Array(9).fill("9"), OP.multiply, ...Array(9).fill("9"), "=");
  await expect(display(page)).toContainText("e");
});
