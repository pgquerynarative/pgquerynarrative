import { useState, useEffect } from "react";
import { ChevronRight, ChevronDown, Database, Table2, Hash } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { api, type SchemaInfo } from "@/api/client";
import { cn } from "@/lib/utils";

interface SchemaBrowserProps {
  onInsert: (text: string) => void;
  className?: string;
}

export function SchemaBrowser({ onInsert, className }: SchemaBrowserProps) {
  const [schemas, setSchemas] = useState<SchemaInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  useEffect(() => {
    api
      .getSchema()
      .then((r) => setSchemas(r.schemas ?? []))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  const toggle = (key: string) => {
    setExpanded((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  if (loading) {
    return (
      <div
        className={cn(
          "rounded-lg overflow-hidden min-w-[220px] bg-surface/95 backdrop-blur-sm border border-primary/10 shadow-md shadow-black/15",
          className
        )}
      >
        <div className="flex items-center gap-2 px-3 py-2 border-b border-border/80 bg-secondary/30">
          <Skeleton className="h-4 w-4 rounded" />
          <Skeleton className="h-4 w-24" />
        </div>
        <div className="p-3 space-y-2">
          <Skeleton className="h-6 w-full" />
          <Skeleton className="h-5 w-[85%]" />
          <Skeleton className="h-5 w-[75%]" />
          <Skeleton className="h-5 w-[90%]" />
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={cn("rounded-lg border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive", className)}>
        {error}
      </div>
    );
  }

  if (schemas.length === 0) {
    return (
      <div
        className={cn(
          "rounded-lg border border-primary/10 bg-surface/95 backdrop-blur-sm p-4 text-sm text-muted-foreground",
          className
        )}
      >
        No schemas available.
      </div>
    );
  }

  /* Schema panel: glassmorphism-lite, 1px cyan border, soft shadow */
  return (
    <div
      className={cn(
        "rounded-lg overflow-hidden min-w-[220px]",
        "bg-surface/95 backdrop-blur-sm border border-primary/10",
        "shadow-md shadow-black/15 ring-1 ring-inset ring-white/[0.02]",
        className
      )}
    >
      <div className="flex items-center gap-2 px-3 py-2 border-b border-border/80 bg-secondary/30">
        <Database className="h-4 w-4 text-primary" />
        <span className="text-sm font-medium">Schema</span>
      </div>
      <div className="max-h-[320px] overflow-auto p-2">
        {schemas.map((schema) => {
          const schemaKey = `schema:${schema.name}`;
          const isSchemaOpen = expanded[schemaKey] ?? true;
          return (
            <div key={schema.name} className="text-sm">
              <button
                type="button"
                onClick={() => toggle(schemaKey)}
                className="flex items-center gap-1.5 w-full py-1.5 px-2 rounded hover:bg-secondary/50 text-left"
              >
                {isSchemaOpen ? <ChevronDown className="h-3.5 w-3.5 shrink-0" /> : <ChevronRight className="h-3.5 w-3.5 shrink-0" />}
                <span className="font-medium text-foreground">{schema.name}</span>
              </button>
              {isSchemaOpen && (
                <div className="ml-4 pl-2 border-l border-border/50 space-y-0.5">
                  {(schema.tables ?? []).map((tbl) => {
                    const tableKey = `${schemaKey}:${tbl.name}`;
                    const isTableOpen = expanded[tableKey] ?? false;
                    const fullTable = `${schema.name}.${tbl.name}`;
                    return (
                      <div key={tbl.name}>
                        <div className="flex items-center gap-1 w-full py-1 px-2 rounded hover:bg-secondary/50">
                          <button type="button" onClick={() => toggle(tableKey)} className="p-0.5 -m-0.5 shrink-0 text-muted-foreground hover:text-foreground">
                            {isTableOpen ? <ChevronDown className="h-3 w-3" /> : <ChevronRight className="h-3 w-3" />}
                          </button>
                          <button
                            type="button"
                            onClick={() => onInsert(fullTable)}
                            className="flex items-center gap-1.5 flex-1 min-w-0 text-left text-muted-foreground hover:text-foreground"
                            title={`Insert ${fullTable}`}
                          >
                            <Table2 className="h-3 w-3 shrink-0" />
                            <span className="font-mono text-xs truncate">{tbl.name}</span>
                          </button>
                        </div>
                        {isTableOpen && (
                          <div className="ml-6 pl-2 border-l border-border/30 space-y-0.5">
                            <button
                              type="button"
                              onClick={() => onInsert(`SELECT * FROM ${fullTable} LIMIT 10`)}
                              className="flex items-center gap-1.5 w-full py-0.5 px-2 rounded hover:bg-primary/10 text-left text-xs text-primary truncate"
                              title={`SELECT * FROM ${fullTable} LIMIT 10`}
                            >
                              SELECT * … LIMIT 10
                            </button>
                            {(tbl.columns ?? []).map((col) => (
                              <button
                                key={col.name}
                                type="button"
                                onClick={() => onInsert(col.name)}
                                className="flex items-center gap-1.5 w-full py-0.5 px-2 rounded hover:bg-secondary/50 text-left"
                              >
                                <Hash className="h-2.5 w-2.5 shrink-0 text-muted-foreground" />
                                <span className="font-mono text-xs">{col.name}</span>
                                <span className="text-[10px] text-muted-foreground truncate">{col.type}</span>
                              </button>
                            ))}
                          </div>
                        )}
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}
