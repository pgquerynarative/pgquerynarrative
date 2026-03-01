import { useEffect, useState, useCallback } from "react";
import { useParams, Link } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { api, type Report } from "@/api/client";
import { FileText, Download, Clock, Cpu, ArrowLeft } from "lucide-react";
import { truncate } from "@/lib/utils";

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
    </div>
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
          <CardContent className="py-16 text-center">
            <FileText className="h-10 w-10 text-muted-foreground mx-auto mb-3" />
            <p className="text-sm text-muted-foreground">No reports yet. <Link to="/query" className="text-primary hover:underline">Run a query</Link> and click Generate Report.</p>
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
