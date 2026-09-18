import type { ReactElement } from "react";
import { Link } from "react-router";

import { Button } from "@/components/ui/button";

export function NotFoundPage(): ReactElement {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-4 text-center">
      <h1 className="text-3xl font-semibold">Página não encontrada</h1>
      <p className="text-muted-foreground">O endereço acessado não existe.</p>
      <Button asChild>
        <Link to="/">Voltar ao início</Link>
      </Button>
    </div>
  );
}
