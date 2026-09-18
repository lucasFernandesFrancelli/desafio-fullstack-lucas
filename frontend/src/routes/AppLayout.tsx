import type { ReactElement } from "react";
import { NavLink, Outlet, useNavigate } from "react-router";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth } from "@/context/AuthContext";
import { cn } from "@/lib/utils";
import { ROLE_LABELS } from "@/lib/constants";

const navItems: Array<{ to: string; label: string }> = [
  { to: "/solicitations", label: "Lista" },
  { to: "/kanban", label: "Kanban" },
];

const gestorNavItems: Array<{ to: string; label: string }> = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/admin", label: "Gestão" },
];

function navLinkClassName({ isActive }: { isActive: boolean }): string {
  return cn(
    "rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
    isActive
      ? "bg-primary/10 text-primary"
      : "text-muted-foreground hover:bg-muted hover:text-foreground",
  );
}

export function AppLayout(): ReactElement {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout(): void {
    logout();
    navigate("/", { replace: true });
  }

  return (
    <div className="flex min-h-svh flex-col bg-background">
      <div className="h-1 bg-brand-gradient" aria-hidden="true" />
      <header className="border-b border-border">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between gap-4 px-4">
          <div className="flex items-center gap-6">
            <span className="flex items-center gap-2 text-sm font-semibold tracking-tight">
              <span className="flex size-7 items-center justify-center rounded-lg bg-brand-gradient text-xs font-bold text-white shadow-sm">
                eK
              </span>
              eKaizen · Melhoria Contínua
            </span>
            <nav className="flex items-center gap-1">
              {navItems.map((item) => (
                <NavLink key={item.to} to={item.to} className={navLinkClassName}>
                  {item.label}
                </NavLink>
              ))}
              {user?.role === "gestor" &&
                gestorNavItems.map((item) => (
                  <NavLink key={item.to} to={item.to} className={navLinkClassName}>
                    {item.label}
                  </NavLink>
                ))}
            </nav>
          </div>

          <div className="flex items-center gap-3">
            <Button size="sm" onClick={() => navigate("/solicitations/new")}>
              Nova solicitação
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="gap-2">
                  <span className="text-sm font-medium">{user?.name}</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>
                  <div className="flex flex-col">
                    <span>{user?.name}</span>
                    <span className="text-xs font-normal text-muted-foreground">
                      {user ? ROLE_LABELS[user.role] : ""}
                    </span>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem onSelect={handleLogout}>Trocar de perfil</DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </header>

      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-6">
        <Outlet />
      </main>
    </div>
  );
}
