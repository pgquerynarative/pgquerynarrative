import { createContext, useCallback, useContext, useState, type ReactNode } from "react";

type Politeness = "polite" | "assertive";

interface AnnounceContextValue {
  announce: (message: string, politeness?: Politeness) => void;
}

const AnnounceContext = createContext<AnnounceContextValue | null>(null);

export function AnnounceProvider({ children }: { children: ReactNode }) {
  const [message, setMessage] = useState("");
  const [politeness, setPoliteness] = useState<Politeness>("polite");

  const announce = useCallback((msg: string, pol: Politeness = "polite") => {
    setPoliteness(pol);
    setMessage(msg);
    // Clear so the same message can be re-announced if triggered again
    const t = setTimeout(() => setMessage(""), 100);
    return () => clearTimeout(t);
  }, []);

  return (
    <AnnounceContext.Provider value={{ announce }}>
      {children}
      {/* Screen reader announcements: aria-live region so messages are spoken */}
      <div
        id="announce-region"
        role="status"
        aria-live={politeness}
        aria-atomic
        className="sr-only"
      >
        {message}
      </div>
    </AnnounceContext.Provider>
  );
}

export function useAnnounce(): AnnounceContextValue["announce"] {
  const ctx = useContext(AnnounceContext);
  if (!ctx) return () => {};
  return ctx.announce;
}
