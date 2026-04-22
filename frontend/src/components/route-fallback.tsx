import { Skeleton } from "@/components/ui/skeleton";
import { Loader2 } from "lucide-react";

/**
 * Shown while a lazy route chunk is loading (Suspense fallback).
 * Full-page loading skeleton with accent and "Loading..." so users see progress.
 */
export function RouteFallback() {
  return (
    <div className="space-y-6 panel-accent-top" aria-busy="true" aria-label="Loading">
      <div className="flex items-center gap-3">
        <Loader2 className="h-6 w-6 animate-spin text-primary" aria-hidden />
        <span className="text-sm font-medium text-muted-foreground">Loading...</span>
      </div>
      <div className="space-y-2">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-4 w-72" />
      </div>
      <div className="grid gap-4 md:grid-cols-3">
        <Skeleton className="h-28 rounded-lg" />
        <Skeleton className="h-28 rounded-lg" />
        <Skeleton className="h-28 rounded-lg" />
      </div>
      <div className="grid gap-4 md:grid-cols-2">
        <Skeleton className="h-48 rounded-lg" />
        <Skeleton className="h-48 rounded-lg" />
      </div>
    </div>
  );
}
