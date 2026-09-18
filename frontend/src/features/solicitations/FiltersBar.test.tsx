import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { FiltersBar } from "@/features/solicitations/FiltersBar";
import { categories } from "@/tests/fixtures";

describe("FiltersBar", () => {
  it("chama onChange com o termo de busca após o debounce", async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();

    render(<FiltersBar filter={{}} onChange={onChange} categories={categories} />);

    await user.type(screen.getByPlaceholderText(/Buscar por título/), "piso");

    expect(onChange).not.toHaveBeenCalled();

    await waitFor(() => expect(onChange).toHaveBeenCalledWith({ q: "piso" }), { timeout: 1000 });
  });

  it("mostra as categorias recebidas via props", () => {
    render(<FiltersBar filter={{}} onChange={vi.fn()} categories={categories} />);
    expect(screen.getByText("Todas as categorias")).toBeInTheDocument();
  });
});
