import type { ReactElement } from "react";
import { Controller, useForm } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import { analysisFinalizeSchema } from "@/schemas/analysisSchema";
import type { AnalysisInput } from "@/types/api";

const SCORE_OPTIONS = [1, 2, 3, 4, 5] as const;

interface AnalysisFormValues {
  severity: string;
  urgency: string;
  trend: string;
  analysisNotes: string;
}

interface AnalysisPanelProps {
  defaultValues: {
    severity: number | null;
    urgency: number | null;
    trend: number | null;
    analysisNotes: string;
  };
  onSavePartial: (input: AnalysisInput) => Promise<void>;
  onFinalize: (input: AnalysisInput) => Promise<void>;
  isSavingPartial: boolean;
  isFinalizing: boolean;
}

function toFormValues(values: AnalysisPanelProps["defaultValues"]): AnalysisFormValues {
  return {
    severity: values.severity !== null ? String(values.severity) : "",
    urgency: values.urgency !== null ? String(values.urgency) : "",
    trend: values.trend !== null ? String(values.trend) : "",
    analysisNotes: values.analysisNotes,
  };
}

export function AnalysisPanel({
  defaultValues,
  onSavePartial,
  onFinalize,
  isSavingPartial,
  isFinalizing,
}: AnalysisPanelProps): ReactElement {
  const form = useForm<AnalysisFormValues>({ defaultValues: toFormValues(defaultValues) });
  const busy = isSavingPartial || isFinalizing;

  function currentPartialInput(values: AnalysisFormValues): AnalysisInput {
    const input: AnalysisInput = { analysisNotes: values.analysisNotes };
    if (values.severity) input.severity = Number(values.severity);
    if (values.urgency) input.urgency = Number(values.urgency);
    if (values.trend) input.trend = Number(values.trend);
    return input;
  }

  async function handleSavePartial(): Promise<void> {
    await onSavePartial(currentPartialInput(form.getValues()));
  }

  async function handleFinalize(): Promise<void> {
    const values = form.getValues();
    const result = analysisFinalizeSchema.safeParse(values);

    if (!result.success) {
      for (const issue of result.error.issues) {
        const field = issue.path[0];
        if (typeof field === "string" && field in values) {
          form.setError(field as keyof AnalysisFormValues, { message: issue.message });
        }
      }
      return;
    }

    await onFinalize(result.data);
  }

  return (
    <div className="flex flex-col gap-4 rounded-lg border border-border p-4">
      <p className="text-sm font-medium">Registre o parecer e as notas para finalizar a análise.</p>

      <div className="grid gap-4 sm:grid-cols-3">
        {(["severity", "urgency", "trend"] as const).map((field) => (
          <div key={field} className="flex flex-col gap-1.5">
            <Label htmlFor={field}>
              {field === "severity" && "Gravidade"}
              {field === "urgency" && "Urgência"}
              {field === "trend" && "Tendência"}
            </Label>
            <Controller
              control={form.control}
              name={field}
              render={({ field: controllerField }) => (
                <Select value={controllerField.value} onValueChange={controllerField.onChange}>
                  <SelectTrigger id={field} className="w-full">
                    <SelectValue placeholder="1 a 5" />
                  </SelectTrigger>
                  <SelectContent>
                    {SCORE_OPTIONS.map((score) => (
                      <SelectItem key={score} value={String(score)}>
                        {score}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              )}
            />
            {form.formState.errors[field] && (
              <p className="text-xs text-destructive">{form.formState.errors[field]?.message}</p>
            )}
          </div>
        ))}
      </div>

      <div className="flex flex-col gap-1.5">
        <Label htmlFor="analysisNotes">Parecer</Label>
        <Textarea id="analysisNotes" rows={4} {...form.register("analysisNotes")} />
        {form.formState.errors.analysisNotes && (
          <p className="text-xs text-destructive">{form.formState.errors.analysisNotes.message}</p>
        )}
      </div>

      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={() => void handleSavePartial()} disabled={busy}>
          {isSavingPartial ? "Salvando…" : "Salvar parecer"}
        </Button>
        <Button onClick={() => void handleFinalize()} disabled={busy}>
          {isFinalizing ? "Finalizando…" : "Finalizar"}
        </Button>
      </div>
    </div>
  );
}
