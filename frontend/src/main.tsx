import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import "@/styles/globals.css";
import { ThemeProvider } from "@/contexts/theme-context";
import Layout from "@/components/layout";
import Dashboard from "@/pages/dashboard";
import QueryRunner from "@/pages/query-runner";
import SavedQueries from "@/pages/saved-queries";
import Reports from "@/pages/reports";
import SettingsPage from "@/pages/settings";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route index element={<Dashboard />} />
            <Route path="query" element={<QueryRunner />} />
            <Route path="saved" element={<SavedQueries />} />
            <Route path="reports" element={<Reports />} />
            <Route path="reports/:id" element={<Reports />} />
            <Route path="settings" element={<SettingsPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  </StrictMode>
);
