import { afterEach, describe, expect, it } from "vitest";
import {
  DEMO_GUIDE_DISMISSED_KEY,
  DEMO_GUIDE_SESSION_KEY,
  dismissDemoGuidePermanently,
  isDemoGuideDismissed,
  markDemoGuidePending,
  shouldShowDemoGuide,
} from "./demo-guide";

describe("demo-guide storage", () => {
  afterEach(() => {
    sessionStorage.clear();
    localStorage.clear();
  });

  it("marks and shows a pending guided-demo landing", () => {
    expect(shouldShowDemoGuide()).toBe(false);
    markDemoGuidePending();
    expect(sessionStorage.getItem(DEMO_GUIDE_SESSION_KEY)).toBe("1");
    expect(shouldShowDemoGuide()).toBe(true);
  });

  it("does not mark pending after permanent dismiss", () => {
    dismissDemoGuidePermanently();
    expect(isDemoGuideDismissed()).toBe(true);
    markDemoGuidePending();
    expect(sessionStorage.getItem(DEMO_GUIDE_SESSION_KEY)).toBeNull();
    expect(shouldShowDemoGuide()).toBe(false);
  });

  it("clears the session flag on dismiss", () => {
    markDemoGuidePending();
    dismissDemoGuidePermanently();
    expect(localStorage.getItem(DEMO_GUIDE_DISMISSED_KEY)).toBe("1");
    expect(sessionStorage.getItem(DEMO_GUIDE_SESSION_KEY)).toBeNull();
    expect(shouldShowDemoGuide()).toBe(false);
  });
});
