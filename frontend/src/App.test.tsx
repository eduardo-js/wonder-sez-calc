import { render, screen } from "@testing-library/react";
import App from "./App";

describe("App", () => {
  it("mounts the calculator screen", () => {
    render(<App />);
    expect(screen.getByTestId("calculator-keypad")).toBeInTheDocument();
  });
});
