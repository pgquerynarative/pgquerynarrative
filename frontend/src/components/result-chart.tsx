import { useState, useEffect, useRef } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  Legend,
} from "recharts";
import { formatFloat } from "@/lib/utils";

export type ChartType = "bar" | "line" | "pie" | "table";

export interface ChartDataPoint {
  name: string;
  [key: string]: string | number;
}

// Finds the index of the first column that looks numeric (number type or numeric values in rows).
function findNumericColumnIndex(columns: { name: string; type: string }[], rows: unknown[][]): number {
  for (let c = 0; c < columns.length; c++) {
    const type = (columns[c].type || "").toLowerCase();
    if (type.includes("int") || type.includes("numeric") || type.includes("decimal") || type.includes("float") || type.includes("double") || type.includes("real")) {
      return c;
    }
    // Check first row value
    if (rows.length > 0 && rows[0][c] != null && typeof rows[0][c] === "number") return c;
  }
  return -1;
}

export function buildChartData(
  columns: { name: string; type: string }[],
  rows: unknown[][]
): ChartDataPoint[] {
  if (!columns.length || !rows.length) return [];
  const labelCol = 0;
  const numCol = findNumericColumnIndex(columns, rows);
  const numericCols = numCol >= 0 ? [numCol] : [];
  if (numericCols.length === 0) {
    for (let c = 1; c < columns.length; c++) {
      if (rows[0][c] != null && (typeof rows[0][c] === "number" || (typeof rows[0][c] === "string" && !isNaN(Number(rows[0][c]))))) {
        numericCols.push(c);
      }
    }
  }
  if (numericCols.length === 0) return [];

  return rows.slice(0, 50).map((row) => {
    const name = row[labelCol] == null ? "" : String(row[labelCol]);
    const out: ChartDataPoint = { name };
    numericCols.forEach((ci) => {
      const v = row[ci];
      const num = typeof v === "number" ? v : Number(v);
      out[columns[ci].name] = isNaN(num) ? 0 : num;
    });
    return out;
  });
}

// All numeric column names for multi-series (e.g. line chart).
export function getNumericColumnNames(
  columns: { name: string; type: string }[],
  rows: unknown[][]
): string[] {
  const names: string[] = [];
  for (let c = 0; c < columns.length; c++) {
    const type = (columns[c].type || "").toLowerCase();
    const isNumericType = type.includes("int") || type.includes("numeric") || type.includes("decimal") || type.includes("float") || type.includes("double") || type.includes("real");
    const isNumericValue = rows.length > 0 && rows[0][c] != null && (typeof rows[0][c] === "number" || (typeof rows[0][c] === "string" && !isNaN(Number(rows[0][c]))));
    if (isNumericType || isNumericValue) names.push(columns[c].name);
  }
  return names;
}

// Vibrant palette on white background
const CHART_COLORS = [
  "#2563eb", // blue
  "#ea580c", // orange
  "#16a34a", // green
  "#dc2626", // red
  "#7c3aed", // violet
  "#0d9488", // teal
  "#ca8a04", // amber
  "#db2777", // pink
];

const CHART_WIDTH = 600;
const CHART_HEIGHT = 320;

interface ResultChartProps {
  chartType: ChartType;
  columns: { name: string; type: string }[];
  rows: unknown[][];
  className?: string;
}

export function ResultChart({ chartType, columns, rows, className = "" }: ResultChartProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [dimensions, setDimensions] = useState({ width: CHART_WIDTH, height: CHART_HEIGHT });

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const ro = new ResizeObserver((entries) => {
      const { width, height } = entries[0]?.contentRect ?? {};
      if (width > 0 && height > 0) setDimensions({ width, height });
    });
    ro.observe(el);
    const { width, height } = el.getBoundingClientRect();
    if (width > 0 && height > 0) setDimensions({ width, height });
    return () => ro.disconnect();
  }, [chartType]);

  const data = buildChartData(columns, rows);
  const numericNames = getNumericColumnNames(columns, rows);
  const valueKey = numericNames[0];

  if (!data.length || !valueKey) {
    return (
      <div className={`flex items-center justify-center rounded-lg border border-border bg-muted/30 p-8 text-sm text-muted-foreground ${className}`}>
        Not enough data to render this chart (need at least one numeric column).
      </div>
    );
  }

  const axisStyle = { tick: { fill: "#374151", fontSize: 11 }, axisLine: { stroke: "#374151" } };
  const gridStroke = "#e5e7eb";
  const tooltipStyle = { backgroundColor: "#fff", border: "1px solid #e5e7eb", borderRadius: "8px", boxShadow: "0 1px 3px rgba(0,0,0,0.1)" };
  const formatter = (v: unknown) => [v != null && typeof v === "number" ? formatFloat(v) : String(v ?? ""), valueKey];

  const chartWrapper = `rounded-lg border border-gray-200 bg-white p-4 overflow-auto ${className}`;

  if (chartType === "bar") {
    return (
      <div ref={containerRef} className={chartWrapper} style={{ minHeight: CHART_HEIGHT, minWidth: 300 }}>
        <BarChart width={dimensions.width} height={CHART_HEIGHT} data={data} margin={{ top: 12, right: 12, left: 0, bottom: 24 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={gridStroke} />
          <XAxis dataKey="name" tick={axisStyle.tick} axisLine={axisStyle.axisLine} />
          <YAxis tick={axisStyle.tick} axisLine={axisStyle.axisLine} tickFormatter={(v) => formatFloat(v)} />
          <Tooltip formatter={formatter} contentStyle={tooltipStyle} />
          <Bar dataKey={valueKey} fill={CHART_COLORS[0]} radius={[4, 4, 0, 0]} name={valueKey} />
        </BarChart>
      </div>
    );
  }

  if (chartType === "line") {
    const series = numericNames.slice(0, 5);
    return (
      <div ref={containerRef} className={chartWrapper} style={{ minHeight: CHART_HEIGHT, minWidth: 300 }}>
        <LineChart width={dimensions.width} height={CHART_HEIGHT} data={data} margin={{ top: 12, right: 12, left: 0, bottom: 24 }}>
          <CartesianGrid strokeDasharray="3 3" stroke={gridStroke} />
          <XAxis dataKey="name" tick={axisStyle.tick} axisLine={axisStyle.axisLine} />
          <YAxis tick={axisStyle.tick} axisLine={axisStyle.axisLine} tickFormatter={(v) => formatFloat(v)} />
          <Tooltip formatter={formatter} contentStyle={tooltipStyle} />
          <Legend wrapperStyle={{ color: "#374151" }} />
          {series.map((key, i) => (
            <Line key={key} type="monotone" dataKey={key} stroke={CHART_COLORS[i % CHART_COLORS.length]} strokeWidth={2.5} dot={{ r: 4, fill: CHART_COLORS[i % CHART_COLORS.length] }} name={key} />
          ))}
        </LineChart>
      </div>
    );
  }

  if (chartType === "pie") {
    return (
      <div ref={containerRef} className={chartWrapper} style={{ minHeight: CHART_HEIGHT, minWidth: 300 }}>
        <PieChart width={dimensions.width} height={CHART_HEIGHT}>
          <Pie
            data={data}
            dataKey={valueKey}
            nameKey="name"
            cx="50%"
            cy="50%"
            outerRadius={100}
            label={({ name, percent }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}
          >
            {data.map((_, i) => (
              <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} stroke="#fff" strokeWidth={2} />
            ))}
          </Pie>
          <Tooltip formatter={formatter} contentStyle={tooltipStyle} />
        </PieChart>
      </div>
    );
  }

  return null;
}
