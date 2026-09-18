import type { ReactElement } from "react";
import { Route, Routes } from "react-router";

import { AppLayout } from "@/routes/AppLayout";
import { DashboardPage } from "@/routes/DashboardPage";
import { KanbanPage } from "@/routes/KanbanPage";
import { ListPage } from "@/routes/ListPage";
import { NewSolicitationPage } from "@/routes/NewSolicitationPage";
import { NotFoundPage } from "@/routes/NotFoundPage";
import { ProfileSelectPage } from "@/routes/ProfileSelectPage";
import { RequireAuth } from "@/routes/RequireAuth";
import { RequireGestor } from "@/routes/RequireGestor";
import { SolicitationDetailPage } from "@/routes/SolicitationDetailPage";

export default function App(): ReactElement {
  return (
    <Routes>
      <Route path="/" element={<ProfileSelectPage />} />

      <Route
        element={
          <RequireAuth>
            <AppLayout />
          </RequireAuth>
        }
      >
        <Route path="/solicitations" element={<ListPage />} />
        <Route path="/solicitations/new" element={<NewSolicitationPage />} />
        <Route path="/solicitations/:id" element={<SolicitationDetailPage />} />
        <Route path="/kanban" element={<KanbanPage />} />
        <Route
          path="/dashboard"
          element={
            <RequireGestor>
              <DashboardPage />
            </RequireGestor>
          }
        />
      </Route>

      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}
