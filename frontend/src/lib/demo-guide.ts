/** Session flag: show the post-landing “What to look at” guide after a guided demo. */
export const DEMO_GUIDE_SESSION_KEY = "pgqn_demo_guide";

/** Local preference: user dismissed the guide and does not want it again. */
export const DEMO_GUIDE_DISMISSED_KEY = "pgqn_demo_guide_dismissed";

export function markDemoGuidePending(): void {
  try {
    if (typeof sessionStorage === "undefined") return;
    if (isDemoGuideDismissed()) return;
    sessionStorage.setItem(DEMO_GUIDE_SESSION_KEY, "1");
  } catch {
    // private browsing / disabled storage
  }
}

export function isDemoGuideDismissed(): boolean {
  try {
    return typeof localStorage !== "undefined" && localStorage.getItem(DEMO_GUIDE_DISMISSED_KEY) === "1";
  } catch {
    return false;
  }
}

/** True when a guided demo just finished and the visitor has not dismissed the guide. */
export function shouldShowDemoGuide(): boolean {
  try {
    if (isDemoGuideDismissed()) return false;
    if (typeof sessionStorage === "undefined") return false;
    return sessionStorage.getItem(DEMO_GUIDE_SESSION_KEY) === "1";
  } catch {
    return false;
  }
}

export function dismissDemoGuidePermanently(): void {
  try {
    if (typeof localStorage !== "undefined") {
      localStorage.setItem(DEMO_GUIDE_DISMISSED_KEY, "1");
    }
    if (typeof sessionStorage !== "undefined") {
      sessionStorage.removeItem(DEMO_GUIDE_SESSION_KEY);
    }
  } catch {
    // ignore
  }
}
