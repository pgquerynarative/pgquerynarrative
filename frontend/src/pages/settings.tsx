import { Card, CardHeader, CardTitle, CardContent, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Database, Cpu, Settings2 } from "lucide-react";

export default function SettingsPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
        <p className="text-muted-foreground mt-1">Server configuration (read-only). Change via environment variables.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center gap-3">
            <Database className="h-5 w-5 text-brand-blue" />
            <div>
              <CardTitle>Database</CardTitle>
              <CardDescription>PostgreSQL connection</CardDescription>
            </div>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <Row label="Host" value={envOrDefault("DB_HOST", "localhost")} />
            <Row label="Port" value={envOrDefault("DB_PORT", "5432")} />
            <Row label="Database" value={envOrDefault("DB_NAME", "pgquerynarrative")} />
            <Row label="Read-only user" value={envOrDefault("DB_READONLY_USER", "pgquerynarrative_readonly")} />
            <Row label="SSL mode" value={envOrDefault("DB_SSL_MODE", "disable")} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center gap-3">
            <Cpu className="h-5 w-5 text-brand-indigo" />
            <div>
              <CardTitle>LLM Provider</CardTitle>
              <CardDescription>Narrative generation model</CardDescription>
            </div>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <Row label="Provider" value={envOrDefault("LLM_PROVIDER", "ollama")} />
            <Row label="Model" value={envOrDefault("LLM_MODEL", "llama3")} />
            <Row label="Base URL" value={envOrDefault("LLM_BASE_URL", "http://localhost:11434")} />
            <Row label="API Key" value={envOrDefault("LLM_API_KEY", "—")} masked />
          </CardContent>
        </Card>

        <Card className="md:col-span-2">
          <CardHeader className="flex flex-row items-center gap-3">
            <Settings2 className="h-5 w-5 text-muted-foreground" />
            <div>
              <CardTitle>Application</CardTitle>
              <CardDescription>General settings</CardDescription>
            </div>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <Row label="Allowed schemas" value="demo" />
            <Row label="Max query length" value="10,000 chars" />
            <Row label="Max rows per query" value="1,000" />
            <Row label="Server port" value={envOrDefault("PORT", "8080")} />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function Row({ label, value, masked }: { label: string; value: string; masked?: boolean }) {
  return (
    <div className="flex items-center justify-between py-1.5 border-b border-border/50 last:border-0">
      <span className="text-muted-foreground">{label}</span>
      {masked && value !== "—" ? (
        <Badge variant="secondary">••••••</Badge>
      ) : (
        <span className="font-mono text-xs">{value}</span>
      )}
    </div>
  );
}

function envOrDefault(key: string, fallback: string): string {
  return fallback;
}
