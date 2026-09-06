# Security policy

PgQueryNarrative connects to production PostgreSQL databases and executes SQL against them.
That makes its security boundary the product, not a feature of it — so this policy states what
the boundary actually is, and how to tell us when it does not hold.

## Reporting a vulnerability

**Do not open a public issue for a security problem.**

Report it privately through GitHub:
[**Report a vulnerability**](https://github.com/pgquery-narrative/pgquerynarrative/security/advisories/new).
That opens a private advisory visible only to maintainers.

Please include the version or commit, your configuration (with secrets removed), what you
expected the boundary to do, and what it did instead. A failing query, request, or migration is
worth more than a description.

What to expect:

| | |
|---|---|
| Acknowledgement | within 3 working days |
| Initial assessment | within 10 working days |
| Fix or documented mitigation | tracked in the advisory until closed |

This is a small project, not a vendor with a 24/7 rota. If something is being actively
exploited, say so in the first line and we will prioritise accordingly.

## Supported versions

Fixes land on `main` and ship in the next release. Only the latest minor release receives
security fixes; there are no long-term support branches.

| Version | Supported |
|---------|-----------|
| 2.1.x   | ✅ |
| 2.0.x   | ❌ — upgrade to 2.1.x |
| 1.x     | ❌ |

## The security model

These are the guarantees the project intends to keep. A reproducible break in any of them is a
vulnerability, and worth reporting.

- **Analytical queries are read-only, enforced by database privilege** — not by an application
  flag. The role the tool queries with cannot `INSERT`, `UPDATE`, `DELETE`, or run DDL, and
  that holds even when `transaction_read_only` is lifted (which index cost projection requires
  for hypopg). `tools/db/verify_security.sh` asserts this, and CI runs it on every PR.
- **Only `SELECT`/`WITH` reaches the database**, validated by walking the `pg_query` parse tree
  rather than matching strings. Multiple statements, DDL, and disallowed schemas are rejected.
- **Bind values are never spliced into SQL as syntax.** A value that only looks like a
  timestamp is quoted and escaped into an inert literal, or refused.
- **Tenant isolation is enforced by row-level security**, with the organisation scope set per
  connection. Cross-organisation reads and writes are blocked in the database.
- **Query results are not sent to an external LLM unless explicitly configured.** The core
  investigation loop runs with no model at all; narratives are optional and can use a local one.
- **Nothing is created or altered without a human action.** The tool proposes index DDL and
  rewrites; it never applies them.

## Known limits — not vulnerabilities

- **`APP_ENV=demo` fabricates workspace KPIs** so the demo has something to show against an
  empty database. It is off by default and gated in one place (`app/service/workspace.go`).
  Never enable it on an instance whose numbers anyone will act on.
- **Reports may contain query results** — including whatever your SQL selected. Treat a
  generated report as being as sensitive as the data behind it, particularly share links.
- **The `app` schema holds SQL text.** With a data encryption key configured it is sealed at
  rest; without one it is stored in plaintext.
- **Anything under `app/*` is internal** and outside the `pkg/narrative` SemVer guarantee. A
  breaking change there is not a security issue.

## Scope

In scope: this repository, the published container image, and the Helm/Kubernetes manifests
under `deploy/`.

Out of scope: vulnerabilities in PostgreSQL itself, in a model provider you configure, or
issues that require an attacker to already hold database superuser or host root.
