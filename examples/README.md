# PgQueryNarrative Examples

This directory contains real-world examples demonstrating how to use PgQueryNarrative to solve business problems.

## Available Examples

### 1. Business Scenario: Sales Performance Analysis
**Location:** `business-scenario/`

A comprehensive example showing how a CEO uses PgQueryNarrative to prepare for a board meeting by analyzing sales data across multiple dimensions.

**What it demonstrates:**
- Multi-dimensional data analysis
- Automatic metrics calculation
- Business insight generation
- Real-world problem solving

**Quick Start:**
```bash
cd business-scenario
bash run-business-scenario.sh
```

See [business-scenario/README.md](business-scenario/README.md) for detailed guide.

---

## How to Use Examples

### Prerequisites

1. **Server running:** `make start-docker` or `make start-local` (see [Getting started](../docs/getting-started/quickstart.md)).

2. **Database ready:**
   - Demo data should be seeded
   - Check with: `curl http://localhost:8080/api/v1/queries/saved`

### Running Examples

Each example includes:
- **README.md** - Detailed user guide
- **Scripts** - Automated execution (if available)
- **Documentation** - Scenario description and expected results

### Step-by-Step

1. Navigate to the example directory
2. Read the README.md
3. Follow the user guide
4. Run the scripts or execute manually
5. Review results and insights

---

## Example Structure

Each example follows this structure:

```
example-name/
├── README.md              # User guide and instructions
├── scenario-description.md # Detailed scenario
├── run-script.sh          # Automated execution (if available)
└── results.md             # Expected results and insights
```

---

## Learning Path

**Beginner:**
1. Start with `business-scenario/`
2. Run the automated script
3. Review the results
4. Try modifying queries

**Intermediate:**
1. Execute queries manually via API
2. Generate narrative reports
3. Save and reuse queries
4. Combine multiple analyses

**Advanced:**
1. Create custom scenarios
2. Integrate with your own data
3. Build automated reporting workflows
4. Embed in your applications

---

## Customizing Examples

All examples use the `demo.sales` table. To use your own data:

1. **Update Queries:** Replace `demo.sales` with your schema/table
2. **Modify Filters:** Adjust date ranges, categories, etc.
3. **Add Dimensions:** Include additional analysis dimensions
4. **Create Reports:** Generate narratives for your specific use case

---

## Contributing Examples

To contribute a new example:

1. Create a new directory: `examples/your-example/`
2. Add README.md with user guide
3. Include scenario description
4. Add execution scripts (optional)
5. Document expected results

---

## Quick Links

- [Main README](../README.md)
- [Documentation index](../docs/README.md)
- [API examples](../docs/api/examples.md)
- [Test queries](../test/queries/quick-test-queries.md)
- Web UI: http://localhost:8080/query

---

Refer to the examples to understand PgQueryNarrative usage patterns.
