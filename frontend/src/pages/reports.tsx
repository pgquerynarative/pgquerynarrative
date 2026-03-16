import { useEffect, useState, useCallback } from "react";
import { useParams, Link } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { api, type Report } from "@/api/client";
import { FileText, Download, Clock, Cpu, ArrowLeft, BarChart3 } from "lucide-react";
import { buttonVariants } from "@/components/ui/button";
import { cn, truncate } from "@/lib/utils";

type TimeSeriesEntry = {
  next_period_forecast?: number | null;
  forecast_ci_lower?: number | null;
  forecast_ci_upper?: number | null;
  trend_summary?: { summary?: string } | null;
  anomalies?: Array<{ period_label?: string; value?: number; reason?: string }> | null;
};
type CohortPeriodPoint = { period_label?: string; periodLabel?: string; value?: number };
type CohortEntry = { cohort_label?: string; cohortLabel?: string; periods?: CohortPeriodPoint[]; retention_pct?: number | null; retentionPct?: number | null };
type MetricsPayload = {
  time_series?: Record<string, TimeSeriesEntry> | null;
  correlations?: Array<{ column_a?: string; column_b?: string; pearson?: number; spearman?: number }> | null;
  cohorts?: CohortEntry[] | null;
};

function ReportDetail() {
  const { id } = useParams<{ id: string }>();
  const [report, setReport] = useState<Report | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    api.getReport(id).then(setReport).catch(() => {}).finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="space-y-4"><Skeleton className="h-8 w-64" /><Skeleton className="h-64 w-full" /></div>;
  if (!report) return <p className="text-muted-foreground">Report not found.</p>;

  const { narrative } = report;

  return (
    <div className="space-y-6">
      <Link to="/reports" className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"><ArrowLeft className="h-4 w-4" /> Back to reports</Link>

      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{narrative?.headline || "Report"}</h1>
          <div className="flex items-center gap-3 mt-2 text-xs text-muted-foreground">
            <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{new Date(report.created_at).toLocaleString()}</span>
            <span className="flex items-center gap-1"><Cpu className="h-3 w-3" />{report.llm_provider} / {report.llm_model}</span>
          </div>
        </div>
        <div className="flex gap-2">
          <a href={`/web/reports/export?id=${report.id}`} download><Button variant="outline" size="sm"><Download className="h-4 w-4" /> HTML</Button></a>
          <a href={`/web/reports/export/pdf?id=${report.id}`} download><Button variant="outline" size="sm"><Download className="h-4 w-4" /> PDF</Button></a>
        </div>
      </div>

      <Card>
        <CardHeader><CardTitle className="text-sm text-muted-foreground">SQL Query</CardTitle></CardHeader>
        <CardContent><pre className="p-4 rounded-md bg-background border border-border text-xs font-mono overflow-auto whitespace-pre-wrap">{report.sql}</pre></CardContent>
      </Card>

      {/* Narrative sections */}
      <Card className="border-primary/20">
        <CardContent className="p-6 space-y-5">
          {narrative?.headline && <h2 className="text-xl font-semibold">{narrative.headline}</h2>}

          {narrative?.takeaways && narrative.takeaways.length > 0 && (
            <Section title="Key Takeaways">{narrative.takeaways.map((t, i) => <Li key={i} color="text-primary">{t}</Li>)}</Section>
          )}
          {narrative?.drivers && narrative.drivers.length > 0 && (
            <Section title="Drivers">{narrative.drivers.map((d, i) => <Li key={i} color="text-brand-blue">{d}</Li>)}</Section>
          )}
          {narrative?.limitations && narrative.limitations.length > 0 && (
            <Section title="Limitations">{narrative.limitations.map((l, i) => <Li key={i} color="text-warning">{l}</Li>)}</Section>
          )}
          {narrative?.recommendations && narrative.recommendations.length > 0 && (
            <Section title="Recommendations">{narrative.recommendations.map((r, i) => <Li key={i} color="text-accent">{r}</Li>)}</Section>
          )}
        </CardContent>
      </Card>

      {report.chart_suggestions && report.chart_suggestions.length > 0 && (
        <Card>
          <CardHeader><CardTitle className="text-sm">Suggested Charts</CardTitle></CardHeader>
          <CardContent><div className="flex flex-wrap gap-2">{report.chart_suggestions.map((s, i) => <Badge key={i} variant="outline">{s.label}</Badge>)}</div></CardContent>
        </Card>
      )}

      <ReportMetricsCard metrics={report.metrics as MetricsPayload | undefined} />
    </div>
  );
}

function ReportMetricsCard({ metrics }: { metrics?: MetricsPayload | null }) {
  // Support both snake_case (API) and camelCase (some proxies)
  const m = metrics as Record<string, unknown> | null | undefined;
  if (!m || typeof m !== "object") return null;

  const tsRaw = m.time_series ?? m.timeSeries;
  const corrRaw = m.correlations;
  const cohortsRaw = m.cohorts;
  const ts = tsRaw && typeof tsRaw === "object" && !Array.isArray(tsRaw) && Object.keys(tsRaw as object).length > 0 ? (tsRaw as Record<string, TimeSeriesEntry>) : null;
  const corr = Array.isArray(corrRaw) && corrRaw.length > 0 ? corrRaw as MetricsPayload["correlations"] : null;
  const cohorts = Array.isArray(cohortsRaw) && cohortsRaw.length > 0 ? (cohortsRaw as CohortEntry[]) : null;

  if (!ts && !corr && !cohorts) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center gap-2">
          <BarChart3 className="h-4 w-4 text-muted-foreground" />
          <CardTitle className="text-sm">Analytics</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">No time-series, correlation, or cohort metrics for this report. Run a query with a date column and numeric measures (or cohort + period columns), then generate a new report.</p>
        </CardContent>
      </Card>
    );
  }

  const formatNum = (n: number | null | undefined) => n != null ? n.toLocaleString(undefined, { maximumFractionDigits: 2 }) : "—";
  const val = (o: Record<string, unknown>, a: string, b: string) => o[a] ?? o[b];

  return (
    <Card>
      <CardHeader className="flex flex-row items-center gap-2">
        <BarChart3 className="h-4 w-4 text-muted-foreground" />
        <CardTitle className="text-sm">Analytics</CardTitle>
      </CardHeader>
      <CardContent className="space-y-5">
        {ts && (
          <div className="space-y-4">
            <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide">Time series & forecast</h4>
            {Object.entries(ts).map(([measure, data]) => {
              const d = data as Record<string, unknown>;
              const forecast = val(d, "next_period_forecast", "nextPeriodForecast") as number | null | undefined;
              const ciLower = val(d, "forecast_ci_lower", "forecastCiLower") as number | null | undefined;
              const ciUpper = val(d, "forecast_ci_upper", "forecastCiUpper") as number | null | undefined;
              const trendSummary = (val(d, "trend_summary", "trendSummary") as { summary?: string } | null) ?? null;
              const anomalies = (val(d, "anomalies", "anomalies") as Array<{ period_label?: string; value?: number; reason?: string }> | null) ?? null;
              return (
                <div key={measure} className="rounded-md border border-border bg-muted/30 p-3 space-y-1.5 text-sm">
                  <p className="font-medium">{measure}</p>
                  {forecast != null && (
                    <p className="text-muted-foreground">Forecast (next period): {formatNum(forecast)}</p>
                  )}
                  {ciLower != null && ciUpper != null && (
                    <p className="text-muted-foreground text-xs">Interval: {formatNum(ciLower)} – {formatNum(ciUpper)}</p>
                  )}
                  {trendSummary?.summary && <p className="text-xs">{trendSummary.summary}</p>}
                  {anomalies && anomalies.length > 0 && (
                    <div className="pt-1">
                      <p className="text-xs font-medium text-muted-foreground mb-1">Anomalies</p>
                      <ul className="list-disc list-inside text-xs space-y-0.5">
                        {anomalies.map((a, i) => {
                          const label = String((a as Record<string, unknown>).period_label ?? (a as Record<string, unknown>).periodLabel ?? "—");
                          return <li key={i}>{label}: {formatNum(a.value)} — {a.reason ?? ""}</li>;
                        })}
                      </ul>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
        {corr && (
          <div className="space-y-2">
            <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide">Correlations</h4>
            <div className="overflow-x-auto rounded-md border border-border">
              <table className="w-full text-sm">
                <thead><tr className="border-b border-border bg-muted/50"><th className="text-left p-2">Column A</th><th className="text-left p-2">Column B</th><th className="text-right p-2">Pearson</th><th className="text-right p-2">Spearman</th></tr></thead>
                <tbody>
                  {corr.map((c, i) => {
                    const row = c as Record<string, unknown>;
                    return (
                      <tr key={i} className={cn("border-b border-border/50 last:border-0 transition-colors", i % 2 === 0 ? "bg-transparent" : "bg-muted/15", "hover:bg-primary/5")}>
                        <td className="p-2 font-mono text-xs">{String(row.column_a ?? row.columnA ?? "—")}</td>
                        <td className="p-2 font-mono text-xs">{String(row.column_b ?? row.columnB ?? "—")}</td>
                        <td className="p-2 text-right">{row.pearson != null ? Number(row.pearson).toFixed(3) : "—"}</td>
                        <td className="p-2 text-right">{row.spearman != null ? Number(row.spearman).toFixed(3) : "—"}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        )}
        {cohorts && (
          <div className="space-y-2">
            <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide">Cohorts</h4>
            <div className="space-y-3">
              {cohorts.map((c, i) => {
                const label = String((c as Record<string, unknown>).cohort_label ?? (c as Record<string, unknown>).cohortLabel ?? "—");
                const periods = Array.isArray(c.periods) ? c.periods : [];
                const ret = (c as Record<string, unknown>).retention_pct ?? (c as Record<string, unknown>).retentionPct;
                return (
                  <div key={i} className="rounded-md border border-border bg-muted/30 p-3 space-y-2 text-sm">
                    <div className="flex items-center justify-between">
                      <p className="font-medium">{label}</p>
                      {ret != null && typeof ret === "number" && (
                        <span className="text-xs text-muted-foreground">Retention: {formatNum(ret)}%</span>
                      )}
                    </div>
                    {periods.length > 0 && (
                      <div className="overflow-x-auto">
                        <table className="w-full min-w-[200px] text-xs">
                          <thead><tr className="border-b border-border bg-muted/50"><th className="text-left p-1.5">Period</th><th className="text-right p-1.5">Value</th></tr></thead>
                          <tbody>
                            {periods.map((p, j) => {
                              const pl = String((p as Record<string, unknown>).period_label ?? (p as Record<string, unknown>).periodLabel ?? j);
                              return (
                                <tr key={j} className={cn("border-b border-border/50 last:border-0 transition-colors", j % 2 === 0 ? "bg-transparent" : "bg-muted/15", "hover:bg-primary/5")}>
                                  <td className="p-1.5">{pl}</td>
                                  <td className="p-1.5 text-right">{formatNum(p.value)}</td>
                                </tr>
                              );
                            })}
                          </tbody>
                        </table>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div>
      <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide mb-2">{title}</h4>
      <ul className="space-y-1.5">{children}</ul>
    </div>
  );
}

function Li({ children, color }: { children: React.ReactNode; color: string }) {
  return <li className="text-sm leading-relaxed flex gap-2"><span className={`${color} mt-1`}>•</span>{children}</li>;
}

function ReportList() {
  const [reports, setReports] = useState<Report[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try { const res = await api.listReports(50, 0); setReports(res.items || []); } catch {}
    finally { setLoading(false); }
  }, []);

  useEffect(() => { load(); }, [load]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Reports</h1>
        <p className="text-muted-foreground mt-1">Generated narrative reports from your queries.</p>
      </div>

      {loading ? (
        <div className="space-y-3">{[1,2,3].map(i => <Skeleton key={i} className="h-20 w-full" />)}</div>
      ) : reports.length === 0 ? (
        <Card>
          <CardContent className="py-16 text-center space-y-4">
            <FileText className="h-10 w-10 text-muted-foreground mx-auto mb-3" />
            <p className="text-sm text-muted-foreground">No reports yet. Run a query and click Generate Report.</p>
            <Link to="/query" className={cn(buttonVariants())}>Run a query</Link>
          </CardContent>
        </Card>
      ) : (
        <div className="space-y-3">
          {reports.map((r) => (
            <Link key={r.id} to={`/reports/${r.id}`}>
              <Card className="hover:border-primary/30 transition-colors cursor-pointer">
                <CardContent className="p-5 flex items-center justify-between">
                  <div className="min-w-0 flex-1">
                    <p className="text-sm font-medium">{r.narrative?.headline || truncate(r.sql, 60)}</p>
                    <div className="flex items-center gap-3 mt-1.5">
                      <span className="text-xs text-muted-foreground flex items-center gap-1"><Clock className="h-3 w-3" />{new Date(r.created_at).toLocaleDateString()}</span>
                      <Badge variant="secondary" className="text-[10px]">{r.llm_provider}</Badge>
                    </div>
                  </div>
                  <FileText className="h-5 w-5 text-brand-cyan flex-shrink-0 ml-4" />
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

export default function Reports() {
  const { id } = useParams();
  return id ? <ReportDetail /> : <ReportList />;
}
