import Calculator from "./components/calculator/Calculator";

/**
 * Root application component. Mounts the Calculator screen.
 */
export default function App() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-100 p-4">
      <Calculator />
    </main>
  );
}
