import { useState, useCallback } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { api, type RunQueryResult, type Report, ApiError } from "@/api/client";
import { Play, FileText, AlertCircle, Clock, Rows3, Download, BarChart3, PieChart as PieChartIcon, LineChart as LineChartIcon, Table2 } from "lucide-react";
import { formatFloat } from "@/lib/utils";
import { ResultChart, type ChartType } from "@/components/result-chart";

export default function QueryRunner() {
  const [sql, setSql] = useState("");
  const [limit, setLimit] = useState("100");
  const [result, setResult] = useState<RunQueryResult | null>(null);
  const [report, setReport] = useState<Report | null>(null);
  const [loading, setLoading] = useState(false);
  const [genLoading, setGenLoading] = useState(false);
  const [error, setError] = useState("");
  const [chartType, setChartType] = useState<ChartType | null>(null);

  const runQuery = useCallback(async () => {
    if (!sql.trim()) { setError("SQL query cannot be empty."); return; }
    setError("");
    setLoading(true);
    setResult(null);
    setReport(null);
    setChartType(null);
    try {
      const r = await api.runQuery(sql, Math.max(1, Math.min(1000, Number(limit) || 100)));
      setResult(r);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Query failed");
    } finally {
      setLoading(false);
    }
  }, [sql, limit]);

  const generateReport = useCallback(async () => {
    if (!sql.trim()) return;
    setGenLoading(true);
    try {
      const r = await api.generateReport(sql);
      setReport(r);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Report generation failed");
    } finally {
      setGenLoading(false);
    }
  }, [sql]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Query Runner</h1>
        <p className="text-muted-foreground mt-1">Execute SQL queries and generate narrative reports.</p>
      </div>

      {/* Editor */}
      <Card>
        <CardContent className="p-6 space-y-4">
          <Textarea
            placeholder="SELECT product_category, SUM(total_amount) AS total&#10;FROM demo.sales&#10;GROUP BY product_category"
            value={sql}
            onChange={(e) => setSql(e.target.value)}
            className="min-h-[160px] font-mono text-sm"
            onKeyDown={(e) => { if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) runQuery(); }}
          />
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div className="flex items-center gap-3">
              <Button onClick={runQuery} disabled={loading}>
                {loading ? <span className="animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full" /> : <Play className="h-4 w-4" />}
                Run Query
              </Button>
              <Button variant="secondary" onClick={generateReport} disabled={genLoading || !sql.trim()}>
                {genLoading ? <span className="animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full" /> : <FileText className="h-4 w-4" />}
                Generate Report
              </Button>
            </div>
            <div className="flex items-center gap-2">
              <label className="text-xs text-muted-foreground whitespace-nowrap">Limit</label>
              <Input type="number" value={limit} onChange={(e) => setLimit(e.target.value)} className="w-20 h-10 text-xs" min={1} max={1000} />
            </div>
          </div>
          <p className="text-[11px] text-muted-foreground">Ctrl+Enter to run. Only SELECT/WITH on the demo schema.</p>
        </CardContent>
      </Card>

      {/* Error */}
      {error && (
        <div className="flex items-start gap-3 p-4 rounded-lg border border-destructive/30 bg-destructive/10 text-destructive text-sm">
          <AlertCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
          <div>{error}</div>
        </div>
      )}

      {/* Loading skeleton */}
      {loading && (
        <Card><CardContent className="p-6 space-y-3">
          <Skeleton className="h-6 w-48" /><Skeleton className="h-48 w-full" />
        </CardContent></Card>
      )}

      {/* Results: chart view or table */}
      {result && !loading && (
        <>
          {/* Suggested charts: click to visualize */}
          {result.chart_suggestions && result.chart_suggestions.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Suggested Charts</CardTitle>
                <p className="text-xs text-muted-foreground mt-1">Click a chart type to visualize the result.</p>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex flex-wrap gap-2">
                  <Button
                    variant={chartType === "table" || chartType === null ? "default" : "outline"}
                    size="sm"
                    onClick={() => setChartType("table")}
                    className="gap-1.5"
                  >
                    <Table2 className="h-3.5 w-3.5" /> Table
                  </Button>
                  {result.chart_suggestions.map((s, i) => {
                    const type = (s.chart_type?.toLowerCase() || "") as ChartType;
                    if (type !== "bar" && type !== "line" && type !== "pie") return null;
                    return (
                      <Button
                        key={i}
                        variant={chartType === type ? "default" : "outline"}
                        size="sm"
                        onClick={() => setChartType(type)}
                        className="gap-1.5"
                        title={s.reason}
                      >
                        {type === "bar" && <BarChart3 className="h-3.5 w-3.5" />}
                        {type === "line" && <LineChartIcon className="h-3.5 w-3.5" />}
                        {type === "pie" && <PieChartIcon className="h-3.5 w-3.5" />}
                        {s.label}
                      </Button>
                    );
                  })}
                </div>
              </CardContent>
            </Card>
          )}

          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Results</CardTitle>
              <div className="flex items-center gap-3 text-xs text-muted-foreground">
                <span className="flex items-center gap-1"><Rows3 className="h-3 w-3" />{result.row_count} rows</span>
                <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{result.execution_time_ms}ms</span>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              {result.rows.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-8">No rows returned.</p>
              ) : (!result.chart_suggestions?.length || chartType === "table" || chartType === null) ? (
                <div className="overflow-auto max-h-[400px]">
                  <table className="w-full text-sm">
                    <thead className="sticky top-0 bg-surface border-b border-border">
                      <tr>
                        {result.columns.map((c, i) => (
                          <th key={i} className="text-left text-xs font-semibold text-muted-foreground px-4 py-3 whitespace-nowrap">{c.name}</th>
                        ))}
                      </tr>
                    </thead>
                    <tbody>
                      {result.rows.map((row, ri) => (
                        <tr key={ri} className="border-b border-border/50 hover:bg-secondary/30 transition-colors">
                          {row.map((cell, ci) => (
                            <td key={ci} className="px-4 py-2.5 whitespace-nowrap font-mono text-xs">{cell == null ? <span className="text-muted-foreground/50">NULL</span> : typeof cell === "number" ? formatFloat(cell) : String(cell)}</td>
                          ))}
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              ) : (
                <div className="p-4">
                  <ResultChart chartType={chartType!} columns={result.columns} rows={result.rows} />
                </div>
              )}
            </CardContent>
          </Card>
        </>
      )}

      {/* Report narrative */}
      {report && (
        <Card className="border-primary/20">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-brand-cyan">Narrative Report</CardTitle>
              <p className="text-xs text-muted-foreground mt-1">{report.llm_provider} / {report.llm_model}</p>
            </div>
            <a href={`/web/reports/export/pdf?id=${report.id}`} download>
              <Button variant="ghost" size="sm"><Download className="h-4 w-4" /> PDF</Button>
            </a>
          </CardHeader>
          <CardContent className="space-y-4">
            {report.narrative.headline && <h3 className="text-lg font-semibold">{report.narrative.headline}</h3>}
            {report.narrative.takeaways && report.narrative.takeaways.length > 0 && (
              <div>
                <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide mb-2">Key Takeaways</h4>
                <ul className="space-y-1.5">{report.narrative.takeaways.map((t, i) => <li key={i} className="text-sm leading-relaxed flex gap-2"><span className="text-primary mt-1">•</span>{t}</li>)}</ul>
              </div>
            )}
            {report.narrative.drivers && report.narrative.drivers.length > 0 && (
              <div>
                <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide mb-2">Drivers</h4>
                <ul className="space-y-1.5">{report.narrative.drivers.map((d, i) => <li key={i} className="text-sm leading-relaxed flex gap-2"><span className="text-brand-blue mt-1">•</span>{d}</li>)}</ul>
              </div>
            )}
            {report.narrative.recommendations && report.narrative.recommendations.length > 0 && (
              <div>
                <h4 className="text-xs font-semibold uppercase text-muted-foreground tracking-wide mb-2">Recommendations</h4>
                <ul className="space-y-1.5">{report.narrative.recommendations.map((r, i) => <li key={i} className="text-sm leading-relaxed flex gap-2"><span className="text-accent mt-1">•</span>{r}</li>)}</ul>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
}
