import { describe, expect, it } from "vitest";

import { formatDate, formatDateTime, formatPriority } from "@/lib/format";

describe("format", () => {
  it("formatDateTime formata data e hora em pt-BR", () => {
    expect(formatDateTime("2026-03-05T14:30:00Z")).toEqual(expect.any(String));
  });

  it("formatDate formata apenas a data em pt-BR", () => {
    expect(formatDate("2026-03-05T14:30:00Z")).toEqual(expect.any(String));
  });

  it("formatPriority retorna travessão quando null", () => {
    expect(formatPriority(null)).toBe("—");
  });

  it("formatPriority retorna o número como string quando presente", () => {
    expect(formatPriority(42)).toBe("42");
  });
});
