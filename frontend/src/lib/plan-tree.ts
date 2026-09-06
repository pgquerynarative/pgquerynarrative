/** Severity levels for plan node visualization. */
export type PlanSeverity = "normal" | "attention" | "high" | "critical";

/** One node in the EXPLAIN plan tree. */
export interface PlanTreeNode {
  id: string;
  nodeType: string;
  label: string;
  schema?: string;
  relation?: string;
  estimatedRows?: number;
  actualRows?: number;
  estimateError?: number;
  cost?: number;
  actualTimeMs?: number;
  loops?: number;
  buffers?: string;
  tempReads?: number;
  tempWrites?: number;
  rowsRemovedByFilter?: number;
  filter?: string;
  severity: PlanSeverity;
  children: PlanTreeNode[];
}

/** Parse PostgreSQL EXPLAIN (FORMAT JSON) into a plan tree. */
export function parsePlanTree(plan: unknown): PlanTreeNode | null {
  if (!plan || typeof plan !== "object") return null;
  const roots = plan as Array<{ Plan?: Record<string, unknown> }>;
  if (!Array.isArray(roots) || roots.length === 0 || !roots[0]?.Plan) return null;
  return walkNode(roots[0].Plan, "0");
}

function walkNode(node: Record<string, unknown>, id: string): PlanTreeNode {
  const nodeType = String(node["Node Type"] ?? "Unknown");
  const schema = node["Schema"] as string | undefined;
  const relation = node["Relation Name"] as string | undefined;
  const estimatedRows = num(node["Plan Rows"]);
  const actualRows = num(node["Actual Rows"]);
  const cost = num(node["Total Cost"]);
  const actualTimeMs = num(node["Actual Total Time"]);
  const loops = num(node["Actual Loops"]) ?? num(node["Plan Loops"]);
  const tempReads = num(node["Temp Read Blocks"]);
  const tempWrites = num(node["Temp Written Blocks"]);
  const rowsRemoved = num(node["Rows Removed by Filter"]);
  const filter = node["Filter"] as string | undefined;

  let estimateError: number | undefined;
  if (estimatedRows && actualRows && estimatedRows > 0) {
    estimateError = actualRows / estimatedRows;
  }

  const severity = computeSeverity(nodeType, estimateError, rowsRemoved, tempWrites);

  let label = nodeType;
  if (relation) {
    label = schema ? `${nodeType}: ${schema}.${relation}` : `${nodeType}: ${relation}`;
  }

  const bufferParts: string[] = [];
  const sharedHit = num(node["Shared Hit Blocks"]);
  const sharedRead = num(node["Shared Read Blocks"]);
  if (sharedHit != null || sharedRead != null) {
    bufferParts.push(`hit=${sharedHit ?? 0} read=${sharedRead ?? 0}`);
  }

  const children: PlanTreeNode[] = [];
  const childPlans = node["Plans"] as unknown[];
  if (Array.isArray(childPlans)) {
    childPlans.forEach((child, i) => {
      if (child && typeof child === "object") {
        children.push(walkNode(child as Record<string, unknown>, `${id}.${i}`));
      }
    });
  }

  return {
    id,
    nodeType,
    label,
    schema,
    relation,
    estimatedRows,
    actualRows,
    estimateError,
    cost,
    actualTimeMs,
    loops,
    buffers: bufferParts.length > 0 ? bufferParts.join(", ") : undefined,
    tempReads,
    tempWrites,
    rowsRemovedByFilter: rowsRemoved,
    filter,
    severity,
    children,
  };
}

function num(v: unknown): number | undefined {
  if (typeof v === "number" && !Number.isNaN(v)) return v;
  return undefined;
}

function computeSeverity(
  nodeType: string,
  estimateError?: number,
  rowsRemoved?: number,
  tempWrites?: number,
): PlanSeverity {
  if (tempWrites && tempWrites > 1000) return "critical";
  if (estimateError && (estimateError >= 100 || estimateError <= 0.01)) return "critical";
  if (nodeType === "Seq Scan" && rowsRemoved && rowsRemoved > 1_000_000) return "high";
  if (estimateError && (estimateError >= 10 || estimateError <= 0.1)) return "high";
  if (nodeType === "Seq Scan") return "attention";
  if (tempWrites && tempWrites > 0) return "attention";
  return "normal";
}

/** Build an evidence finding message for a plan node. */
export function nodeFinding(node: PlanTreeNode): {
  finding: string;
  whyItMatters: string;
  investigate: string[];
  confidence: "high" | "medium" | "low";
} {
  const items: string[] = [];
  let finding = `Plan node: ${node.label}`;
  let why = "This node contributes to overall query cost.";
  let confidence: "high" | "medium" | "low" = "low";

  if (node.estimateError && node.estimatedRows != null && node.actualRows != null) {
    const ratio = node.estimateError;
    finding = `The planner estimated ${formatNum(node.estimatedRows)} rows, but this node returned ${formatNum(node.actualRows)} rows—a ${ratio >= 1 ? `${ratio.toFixed(0)}×` : `${(1 / ratio).toFixed(0)}×`} ${ratio >= 1 ? "underestimation" : "overestimation"}.`;
    why = "The misestimated row count may have affected join strategy and caused significantly more work downstream.";
    confidence = ratio >= 100 || ratio <= 0.01 ? "high" : "medium";
    items.push("Column statistics", "Correlated predicates", "Extended statistics", "Stale ANALYZE data");
  } else if (node.nodeType === "Seq Scan" && node.rowsRemovedByFilter && node.rowsRemovedByFilter > 10000) {
    finding = `Sequential scan on ${node.relation ?? "relation"} removes ${formatNum(node.rowsRemovedByFilter)} rows by filter.`;
    why = "A large fraction of rows are read and discarded, which may indicate a missing or unusable index.";
    confidence = "medium";
    items.push("Index definitions on filter columns", "Predicate shape (functions on indexed columns)", "Table statistics freshness");
  } else if (node.tempWrites && node.tempWrites > 0) {
    finding = `This node wrote ${node.tempWrites} temporary blocks to disk.`;
    why = "Spilling to disk increases latency and I/O pressure.";
    confidence = "high";
    items.push("work_mem setting", "Sort/hash input size", "Row width and aggregation");
  }

  if (items.length === 0) {
    items.push("Predicate selectivity", "Join order", "Parallelism settings");
  }

  return { finding, whyItMatters: why, investigate: items, confidence };
}

function formatNum(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1000) return `${(n / 1000).toFixed(1)}k`;
  return String(Math.round(n));
}

/** Collect all nodes from a plan tree (flat list). */
export function flattenPlanTree(root: PlanTreeNode | null): PlanTreeNode[] {
  if (!root) return [];
  const out: PlanTreeNode[] = [root];
  for (const child of root.children) {
    out.push(...flattenPlanTree(child));
  }
  return out;
}

// Mirrors PARTITION_SUFFIX in ./utils. Duplicated rather than imported because
// utils pulls in clsx/tailwind-merge, and the plan model should not depend on
// presentation helpers.
const PARTITION_SUFFIX = /_\d{4}_\d{2}$/;

const SEVERITY_ORDER: Record<PlanSeverity, number> = {
  normal: 0,
  attention: 1,
  high: 2,
  critical: 3,
};

function sumDefined(values: Array<number | undefined>): number | undefined {
  const present = values.filter((v): v is number => typeof v === "number");
  return present.length > 0 ? present.reduce((a, b) => a + b, 0) : undefined;
}

/**
 * Collapses sibling scans of monthly partitions into one row.
 *
 * A plan over an unpruned partitioned table produces one near-identical child per
 * partition, which buries the rest of the tree. Siblings are grouped when they
 * share a node type, schema and parent relation — that is, their relation names
 * differ only by a trailing `_YYYY_MM`.
 *
 * Additive metrics are summed, so the collapsed row reports the work the whole
 * group did rather than one arbitrary member's share, and severity takes the
 * worst of the group so a critical partition is never hidden behind a quiet one.
 *
 * Two deliberate limits:
 *  - Only childless nodes are grouped. Concatenating the children of N siblings
 *    would invent a subtree shape that no plan actually had.
 *  - A group of one keeps its original label and relation. Rewriting a lone
 *    `sales_2023_01` to `sales` would discard the month while collapsing nothing.
 */
export function collapsePartitionSiblings(nodes: PlanTreeNode[]): PlanTreeNode[] {
  const groups = new Map<string, PlanTreeNode[]>();
  const order: string[] = [];

  for (const node of nodes) {
    const relation = node.relation ?? "";
    const collapsible =
      relation !== "" && PARTITION_SUFFIX.test(relation) && node.children.length === 0;
    // Namespaced apart from pass-through nodes: stripping the suffix from
    // "sales_2023_01" yields "sales", which is also a real relation — a scan of
    // the parent table must not be folded in with its children.
    const key = collapsible
      ? `P|${node.nodeType}|${node.schema ?? ""}|${relation.replace(PARTITION_SUFFIX, "")}`
      : `N|${node.id}`;

    const existing = groups.get(key);
    if (existing) {
      existing.push(node);
    } else {
      groups.set(key, [node]);
      order.push(key);
    }
  }

  return order.map((key) => {
    const group = groups.get(key)!;
    const first = group[0];
    if (group.length === 1) {
      return first;
    }

    const base = (first.relation ?? "").replace(PARTITION_SUFFIX, "");
    const qualified = first.schema ? `${first.schema}.${base}` : base;
    const worst = group.reduce(
      (acc, n) => (SEVERITY_ORDER[n.severity] > SEVERITY_ORDER[acc] ? n.severity : acc),
      first.severity,
    );

    return {
      ...first,
      relation: base,
      label: `${first.nodeType}: ${qualified} (×${group.length} partitions)`,
      severity: worst,
      estimatedRows: sumDefined(group.map((n) => n.estimatedRows)),
      actualRows: sumDefined(group.map((n) => n.actualRows)),
      cost: sumDefined(group.map((n) => n.cost)),
      actualTimeMs: sumDefined(group.map((n) => n.actualTimeMs)),
      tempReads: sumDefined(group.map((n) => n.tempReads)),
      tempWrites: sumDefined(group.map((n) => n.tempWrites)),
      rowsRemovedByFilter: sumDefined(group.map((n) => n.rowsRemovedByFilter)),
      // estimateError is a ratio; summing it would be meaningless, and the group's
      // members can disagree. Left off the collapsed row rather than guessed.
      estimateError: undefined,
      children: [],
    };
  });
}
