import type { ReactElement } from "react";
import { Loader2Icon } from "lucide-react";

interface PageSpinnerProps {
  label?: string;
  /** Usa min-h-svh para ocupar a tela inteira (fora do AppLayout). Dentro do layout, prefira o padrão (só a área de conteúdo). */
  fullScreen?: boolean;
}

export function PageSpinner({ label = "Carregando…", fullScreen = false }: PageSpinnerProps): ReactElement {
  return (
    <div
      className={
        fullScreen
          ? "flex min-h-svh flex-col items-center justify-center gap-3 text-muted-foreground"
          : "flex flex-1 flex-col items-center justify-center gap-3 py-16 text-muted-foreground"
      }
      role="status"
      aria-live="polite"
    >
      <Loader2Icon className="size-6 animate-spin text-primary" aria-hidden="true" />
      <span className="text-sm">{label}</span>
    </div>
  );
}
