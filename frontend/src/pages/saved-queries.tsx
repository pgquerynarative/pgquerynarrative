import { useEffect, useState, useCallback } from "react";
import { Link } from "react-router-dom";
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Button, buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { api, type SavedQuery, ApiError } from "@/api/client";
import { Search, Trash2, Copy, Play, AlertCircle, Bookmark } from "lucide-react";
import { truncate } from "@/lib/utils";
import { useNavigate } from "react-router-dom";

export default function SavedQueries() {
  const [queries, setQueries] = useState<SavedQuery[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [error, setError] = useState("");
  const [selected, setSelected] = useState<SavedQuery | null>(null);
  const nav = useNavigate();

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const res = await api.listSaved(100, 0);
      setQueries(res.items || []);
    } catch { setError("Failed to load saved queries."); }
    finally { setLoading(false); }
  }, []);

  useEffect(() => { load(); }, [load]);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this saved query?")) return;
    try {
      await api.deleteSaved(id);
      setQueries((prev) => prev.filter((q) => q.id !== id));
      if (selected?.id === id) setSelected(null);
    } catch (e) {
      setError(e instanceof ApiError ? e.message : "Delete failed.");
    }
  };

  const filtered = queries.filter(
    (q) => q.name.toLowerCase().includes(search.toLowerCase()) || q.sql.toLowerCase().includes(search.toLowerCase()) || (q.tags || []).some((t) => t.toLowerCase().includes(search.toLowerCase()))
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Saved Queries</h1>
          <p className="text-muted-foreground mt-1">Manage your bookmarked SQL queries.</p>
        </div>
      </div>

      {error && (
        <div
          role="alert"
          className="flex items-start gap-3 p-4 rounded-lg border border-destructive/40 bg-destructive/10 text-destructive text-sm shadow-sm ring-1 ring-destructive/20"
        >
          <AlertCircle className="h-5 w-5 flex-shrink-0" aria-hidden />
          <span className="min-w-0">{error}</span>
        </div>
      )}

      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
        <Input placeholder="Search by name, SQL, or tag..." value={search} onChange={(e) => setSearch(e.target.value)} className="pl-10" />
      </div>

      <div className="grid gap-6 md:grid-cols-[1fr_1fr]">
        {/* List */}
        <div className="space-y-2">
          {loading ? [1,2,3,4].map((i) => <Skeleton key={i} className="h-20 w-full" />) : filtered.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center space-y-4">
                <Bookmark className="h-8 w-8 text-muted-foreground mx-auto mb-3" />
                <p className="text-sm text-muted-foreground">{search ? "No queries match your search." : "No saved queries yet. Save one from the Query Runner."}</p>
                {!search && (
                  <Link to="/query" className={cn(buttonVariants())}>Go to Query Runner</Link>
                )}
              </CardContent>
            </Card>
          ) : filtered.map((q) => (
            <button
              key={q.id}
              onClick={() => setSelected(q)}
              className={`w-full text-left p-4 rounded-lg border transition-colors cursor-pointer ${selected?.id === q.id ? "border-primary/50 bg-primary/5" : "border-border hover:bg-secondary/30"}`}
            >
              <div className="flex items-center justify-between">
                <p className="text-sm font-medium truncate">{q.name}</p>
                <span className="text-[10px] text-muted-foreground">{new Date(q.created_at).toLocaleDateString()}</span>
              </div>
              <p className="text-xs text-muted-foreground font-mono mt-1 truncate">{truncate(q.sql, 80)}</p>
              {q.tags && q.tags.length > 0 && (
                <div className="flex gap-1.5 mt-2">{q.tags.map((t) => <Badge key={t} variant="outline" className="text-[10px]">{t}</Badge>)}</div>
              )}
            </button>
          ))}
        </div>

        {/* Detail panel */}
        <Card className="h-fit sticky top-8">
          {selected ? (
            <>
              <CardHeader>
                <CardTitle>{selected.name}</CardTitle>
                <CardDescription>Created {new Date(selected.created_at).toLocaleString()}</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {selected.description && <p className="text-sm text-muted-foreground">{selected.description}</p>}
                <pre className="p-4 rounded-md bg-background border border-border text-xs font-mono overflow-auto max-h-[200px] whitespace-pre-wrap">{selected.sql}</pre>
                {selected.tags && selected.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1.5">{selected.tags.map((t) => <Badge key={t}>{t}</Badge>)}</div>
                )}
                <div className="flex gap-2 pt-2">
                  <Button size="sm" onClick={() => nav("/query", { state: { sql: selected.sql } })}><Play className="h-3.5 w-3.5" /> Run</Button>
                  <Button variant="ghost" size="sm" onClick={() => navigator.clipboard.writeText(selected.sql)}><Copy className="h-3.5 w-3.5" /> Copy SQL</Button>
                  <Button variant="destructive" size="sm" onClick={() => handleDelete(selected.id)} className="ml-auto"><Trash2 className="h-3.5 w-3.5" /> Delete</Button>
                </div>
              </CardContent>
            </>
          ) : (
            <CardContent className="py-16 text-center">
              <p className="text-sm text-muted-foreground">Select a query to view details.</p>
            </CardContent>
          )}
        </Card>
      </div>
    </div>
  );
}
