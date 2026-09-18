import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { StatusBadge } from "@/features/solicitations/StatusBadge";
import { STATUS_LABELS } from "@/lib/constants";
import type { Status } from "@/types/api";

describe("StatusBadge", () => {
  const statuses: Status[] = ["rascunho", "em_aprovacao", "em_analise", "finalizada", "recusada"];

  it.each(statuses)("mostra o rótulo em português para %s", (status) => {
    render(<StatusBadge status={status} />);
    expect(screen.getByText(STATUS_LABELS[status])).toBeInTheDocument();
  });
});
