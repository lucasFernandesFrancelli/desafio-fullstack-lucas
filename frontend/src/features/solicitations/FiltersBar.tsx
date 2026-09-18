import { useEffect, useState, type ReactElement } from "react";

import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { SolicitationListFilter } from "@/api/solicitations";
import { STATUS_LABELS, STATUS_ORDER } from "@/lib/constants";
import type { Category } from "@/types/api";

const ALL_VALUE = "all";

interface FiltersBarProps {
  filter: SolicitationListFilter;
  onChange: (filter: SolicitationListFilter) => void;
  categories: Category[];
}

export function FiltersBar({ filter, onChange, categories }: FiltersBarProps): ReactElement {
  const [searchTerm, setSearchTerm] = useState(filter.q ?? "");

  useEffect(() => {
    const timeout = setTimeout(() => {
      if (searchTerm !== (filter.q ?? "")) {
        onChange({ ...filter, q: searchTerm || undefined });
      }
    }, 300);
    return () => clearTimeout(timeout);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchTerm]);

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
      <Input
        placeholder="Buscar por título, problema ou melhoria…"
        value={searchTerm}
        onChange={(event) => setSearchTerm(event.target.value)}
        className="sm:max-w-xs"
      />

      <Select
        value={filter.status ?? ALL_VALUE}
        onValueChange={(value) =>
          onChange({ ...filter, status: value === ALL_VALUE ? undefined : (value as SolicitationListFilter["status"]) })
        }
      >
        <SelectTrigger className="sm:w-48">
          <SelectValue placeholder="Estado" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={ALL_VALUE}>Todos os estados</SelectItem>
          {STATUS_ORDER.map((status) => (
            <SelectItem key={status} value={status}>
              {STATUS_LABELS[status]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select
        value={filter.categoryId ?? ALL_VALUE}
        onValueChange={(value) => onChange({ ...filter, categoryId: value === ALL_VALUE ? undefined : value })}
      >
        <SelectTrigger className="sm:w-56">
          <SelectValue placeholder="Categoria" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value={ALL_VALUE}>Todas as categorias</SelectItem>
          {categories.map((category) => (
            <SelectItem key={category.id} value={category.id}>
              {category.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
