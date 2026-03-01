# LLM setup

Report generation requires an LLM. Supported providers: **Ollama** (local), **Gemini**, **Claude**, **OpenAI**, **Groq**. Set [Configuration – LLM](../configuration.md#llm) variables: `LLM_PROVIDER`, `LLM_MODEL`, and for cloud providers `LLM_API_KEY`.

## Ollama (local)

No API key. Runs on your machine.

- **Install:** [ollama.ai](https://ollama.ai) or `brew install ollama` (macOS). Start: `ollama serve` (API: http://localhost:11434).
- **Model:** `ollama pull llama3.2` (default). Others: `mistral`, `llama2`. Verify: `curl http://localhost:11434/api/tags`.
- **Config:** Defaults `LLM_PROVIDER=ollama`, `LLM_MODEL=llama3.2`, `LLM_BASE_URL=http://localhost:11434`. App in Docker, Ollama on host: `LLM_BASE_URL=http://host.docker.internal:11434` (see [Deployment](../reference/deployment.md)).

## Gemini (Google)

- **API key:** [Google AI Studio](https://aistudio.google.com/app/apikey).
- **Config:** `LLM_PROVIDER=gemini`, `LLM_MODEL=gemini-2.0-flash`, `LLM_API_KEY=your_key`. On 404 try another model; on 429 the client retries.

## Claude (Anthropic)

- **API key:** [Anthropic Console](https://console.anthropic.com/). Paid.
- **Config:** `LLM_PROVIDER=claude`, `LLM_MODEL=claude-3-5-sonnet-20241022`, `LLM_API_KEY=your_key`. [Anthropic models](https://docs.anthropic.com/en/docs/about-claude/models).

## OpenAI (GPT)

- **API key:** [OpenAI API keys](https://platform.openai.com/api-keys).
- **Config:** `LLM_PROVIDER=openai`, `LLM_MODEL=gpt-4o-mini`, `LLM_API_KEY=your_key`. Others: `gpt-4o`, `gpt-4-turbo`. On 429 the client retries.

## Groq

- **API key:** [Groq Console](https://console.groq.com/keys). OpenAI-compatible; fast inference.
- **Config:** `LLM_PROVIDER=groq`, `LLM_MODEL=llama-3.3-70b-versatile`, `LLM_API_KEY=your_key`. Others: `llama-3.1-8b-instant`, `mixtral-8x7b-32768`.

## MCP (Claude desktop / Cursor)

Use PgQueryNarrative as tools (run query, generate report, list saved/reports) from Claude or Cursor.

1. Run app and build MCP: `make build-mcp` → `bin/mcp-server`.
2. Edit your client’s MCP config: [Configuration – MCP](../configuration.md#mcp-claude-desktop--cursor).
3. Add the `pgquerynarrative` server; **command** = path to `bin/mcp-server`. Example: `config/mcp-example.json`; copy the `pgquerynarrative` block into `mcpServers`.
4. Restart the client. Tools: `run_query`, `generate_report`, `list_saved_queries`, `get_report`, `list_reports`.

## Troubleshooting

| Issue | Action |
|-------|--------|
| **Connection refused (Ollama)** | Start Ollama (`ollama serve`). Docker: `LLM_BASE_URL=http://host.docker.internal:11434`. Ensure a model is pulled. |
| **Report fails / timeout** | Check provider and API key. [Troubleshooting – Reports and LLM](../reference/troubleshooting.md#reports-and-llm). |
| **Slow (Ollama)** | Use a smaller model (e.g. `mistral`); first run per session is slower. |

## See also

- [Configuration](../configuration.md) — All LLM and MCP variables
- [Quick start](quickstart.md) — Get the app running
- [Deployment](../reference/deployment.md) — Docker/Helm (e.g. LLM_BASE_URL in containers)
- [Embedded integration](embedded.md) — Library and middleware (same config)
- [Troubleshooting](../reference/troubleshooting.md) — Report and connection issues
- [Documentation index](../README.md)
