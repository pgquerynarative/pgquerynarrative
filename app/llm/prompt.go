package llm

import (
	"fmt"
	"strings"
)

// BuildNarrativePrompt creates a prompt for narrative generation from query results
func BuildNarrativePrompt(sql string, columns []string, rows [][]interface{}, metricsJSON string) string {
	var sb strings.Builder

	sb.WriteString("You are a data analyst expert. Your task is to convert SQL query results into a clear, evidence-based business narrative.\n\n")
	sb.WriteString("IMPORTANT RULES:\n")
	sb.WriteString("1. Only make claims that are directly supported by the data provided\n")
	sb.WriteString("2. Cite specific numbers and metrics in your narrative\n")
	sb.WriteString("3. Do not make assumptions or inferences beyond what the data shows\n")
	sb.WriteString("4. Acknowledge limitations if the dataset is small or incomplete\n")
	sb.WriteString("5. Use clear, professional business language\n")
	sb.WriteString("6. If the CALCULATED METRICS include a \"TimeSeries\" (or \"time_series\") object with period-over-period comparison, include at least one takeaway that mentions how key measures changed vs the previous period (e.g. \"Revenue was −0.2% vs the previous period\" or \"Transactions were up 3.1% compared to the prior period\").\n\n")

	sb.WriteString("SQL QUERY:\n")
	sb.WriteString(sql)
	sb.WriteString("\n\n")

	sb.WriteString("QUERY RESULTS:\n")
	sb.WriteString("Columns: ")
	sb.WriteString(strings.Join(columns, ", "))
	sb.WriteString("\n\n")

	// Include sample rows (first 10)
	sampleRows := rows
	if len(sampleRows) > 10 {
		sampleRows = sampleRows[:10]
	}

	sb.WriteString("Sample Data (showing first 10 rows):\n")
	maxCols := 0
	for _, row := range sampleRows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	rowStr := make([]string, maxCols)
	for i, row := range sampleRows {
		sb.WriteString(fmt.Sprintf("Row %d: ", i+1))
		n := len(row)
		for j := 0; j < n; j++ {
			rowStr[j] = fmt.Sprintf("%v", row[j])
		}
		sb.WriteString(strings.Join(rowStr[:n], " | "))
		sb.WriteString("\n")
	}

	if len(rows) > 10 {
		sb.WriteString(fmt.Sprintf("... and %d more rows\n", len(rows)-10))
	}

	sb.WriteString("\n")
	sb.WriteString("CALCULATED METRICS:\n")
	sb.WriteString(metricsJSON)
	sb.WriteString("\n\n")

	sb.WriteString("TASK:\n")
	sb.WriteString("Generate a business narrative in JSON format with the following structure:\n")
	sb.WriteString(`{
  "headline": "A concise one-sentence summary of the key finding",
  "takeaways": [
    "Key insight 1 with specific numbers",
    "Key insight 2 with specific numbers",
    "Key insight 3 with specific numbers"
  ],
  "drivers": [
    "Factor that explains the results (if applicable)"
  ],
  "limitations": [
    "Any limitations or caveats about the data or analysis"
  ],
  "recommendations": [
    "Actionable recommendation based on the data (optional)"
  ]
}`)
	sb.WriteString("\n\n")
	sb.WriteString("Return ONLY valid JSON. Do not include any markdown formatting, code blocks, or explanatory text outside the JSON.\n")

	return sb.String()
}
