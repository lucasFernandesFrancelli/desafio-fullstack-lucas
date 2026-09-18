export function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString("pt-BR", { dateStyle: "short", timeStyle: "short" });
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("pt-BR");
}

export function formatPriority(priority: number | null): string {
  return priority === null ? "—" : String(priority);
}
