const BASE = "/api/v1";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json", ...init?.headers },
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({})) as Record<string, unknown>;
    const message =
      (typeof body?.message === "string" && body.message) ||
      (body?.body && typeof (body.body as Record<string, unknown>).message === "string" && (body.body as Record<string, unknown>).message) ||
      (typeof body?.name === "string" && body.name) ||
      res.statusText;
    throw new ApiError(res.status, message as string, body);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export class ApiError extends Error {
  constructor(public status: number, message: string, public body?: unknown) {
    super(message);
  }
}

export interface Column { name: string; type: string; }
export interface ChartSuggestion { chart_type: string; label: string; reason: string; }
export interface RunQueryResult {
  columns: Column[];
  rows: unknown[][];
  row_count: number;
  execution_time_ms: number;
  limit: number;
  chart_suggestions?: ChartSuggestion[];
}

export interface SavedQuery {
  id: string;
  name: string;
  sql: string;
  description?: string;
  tags?: string[];
  created_at: string;
  updated_at?: string;
}

export interface NarrativeContent {
  headline: string;
  takeaways?: string[];
  drivers?: string[];
  limitations?: string[];
  recommendations?: string[];
}

export interface Report {
  id: string;
  sql: string;
  narrative: NarrativeContent;
  metrics: Record<string, unknown>;
  chart_suggestions?: ChartSuggestion[];
  created_at: string;
  llm_model: string;
  llm_provider: string;
}

export interface AskResult {
  question: string;
  sql: string;
  report: Report;
}

export interface SchemaInfo { name: string; tables: { name: string; columns: Column[] }[]; }

export interface AnalyticsSettings {
  anomaly_sigma: number;
  anomaly_method?: string;
  trend_periods: number;
  moving_avg_window: number;
  trend_threshold_percent: number;
  confidence_level: number;
  min_rows_for_correlation?: number;
  smoothing_alpha?: number;
  smoothing_beta?: number;
  max_seasonal_lag?: number;
  min_periods_for_seasonality?: number;
  max_timeseries_periods?: number;
}

export interface SettingsResponse {
  analytics: AnalyticsSettings;
}

// Normalize SQL for API: trim and strip trailing semicolon (API rejects ";" in sql).
function normalizeSql(sql: string): string {
  return sql.trim().replace(/;\s*$/, "");
}

export const api = {
  runQuery: (sql: string, limit = 100) =>
    request<RunQueryResult>("/queries/run", {
      method: "POST",
      body: JSON.stringify({ sql: normalizeSql(sql), limit }),
    }),

  listSaved: (limit = 50, offset = 0) =>
    request<{ items: SavedQuery[]; limit: number; offset: number }>(`/queries/saved?limit=${limit}&offset=${offset}`),

  saveQuery: (name: string, sql: string, tags: string[] = []) =>
    request<SavedQuery>("/queries/saved", { method: "POST", body: JSON.stringify({ name, sql, tags }) }),

  getSaved: (id: string) => request<SavedQuery>(`/queries/saved/${id}`),

  deleteSaved: (id: string) => request<void>(`/queries/saved/${id}`, { method: "DELETE" }),

  generateReport: (sql: string) =>
    request<Report>("/reports/generate", {
      method: "POST",
      body: JSON.stringify({ sql: normalizeSql(sql) }),
    }),

  listReports: (limit = 50, offset = 0) =>
    request<{ items: Report[]; limit: number; offset: number }>(`/reports?limit=${limit}&offset=${offset}`),

  getReport: (id: string) => request<Report>(`/reports/${id}`),

  getSchema: () => request<{ schemas: SchemaInfo[] }>("/schema"),

  getSettings: () => request<SettingsResponse>("/settings"),

  getSuggestions: (limit = 5) =>
    request<{ suggestions: { sql: string; title: string; source: string }[] }>(`/suggestions/queries?limit=${limit}`),

  ask: (question: string) =>
    request<AskResult>("/suggestions/ask", {
      method: "POST",
      body: JSON.stringify({ question: question.trim() }),
    }),
};
