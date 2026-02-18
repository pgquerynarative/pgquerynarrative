# LLM setup (report generation)

Report generation needs an LLM. Optional; only required for generating narratives. Supported: **Ollama** (local), **Gemini**, **Claude**, **OpenAI (GPT)**, **Groq**. Set `LLM_PROVIDER` and, for cloud providers, `LLM_API_KEY`. Full variable list: [Configuration – LLM](../configuration.md#llm).

## Ollama (local)

No API key. Runs on your machine.

- **Install:** [ollama.ai](https://ollama.ai) or macOS `brew install ollama`; Linux `curl -fsSL https://ollama.ai/install.sh | sh`. Then `ollama serve` (API: `http://localhost:11434`).
- **Model:** `ollama pull llama3.2` (default). Others: `mistral`, `llama2`. Verify: `curl http://localhost:11434/api/tags`.
- **Config:** Defaults `LLM_PROVIDER=ollama`, `LLM_MODEL=llama3.2`, `LLM_BASE_URL=http://localhost:11434`. With app in Docker and Ollama on host: `LLM_BASE_URL=http://host.docker.internal:11434`.

## Gemini (Google)

- **API key:** [Google AI Studio](https://aistudio.google.com/app/apikey).
- **Config:** `LLM_PROVIDER=gemini`, `LLM_MODEL=gemini-2.0-flash` (or `gemini-1.5-flash`), `LLM_API_KEY=your_key`. On 404 try another model; on 429 the client retries.

## Claude (Anthropic)

- **API key:** [Anthropic Console](https://console.anthropic.com/). Paid; no free API tier.
- **Config:** `LLM_PROVIDER=claude`, `LLM_MODEL=claude-3-5-sonnet-20241022`, `LLM_API_KEY=your_key`. See [Anthropic models](https://docs.anthropic.com/en/docs/about-claude/models).

## OpenAI (GPT)

- **API key:** [OpenAI API keys](https://platform.openai.com/api-keys).
- **Config:** `LLM_PROVIDER=openai`, `LLM_MODEL=gpt-4o-mini` (default), `LLM_API_KEY=your_key`. Other models: `gpt-4o`, `gpt-4-turbo`. On 429 the client retries.

## Groq

- **API key:** [Groq Console](https://console.groq.com/keys). OpenAI-compatible API; fast inference.
- **Config:** `LLM_PROVIDER=groq`, `LLM_MODEL=llama-3.3-70b-versatile` (default), `LLM_API_KEY=your_key`. Other models: `llama-3.1-8b-instant`, `mixtral-8x7b-32768`. On 429 the client retries.

## MCP (Claude desktop / Cursor)

Use PgQueryNarrative as tools (run query, generate report, list saved/reports) from Claude app or Cursor.

1. Run app and build MCP: `make build-mcp` → `bin/mcp-server`.
2. Edit your client’s MCP config (see [Configuration – MCP](../configuration.md#mcp-claude-desktop--cursor)).
3. Add the `pgquerynarrative` server; use the **command** path to `bin/mcp-server`. Example config: [config/mcp-example.json](../../config/mcp-example.json). Copy the `pgquerynarrative` block into your `mcpServers`.
4. Restart the client. Tools: `run_query`, `generate_report`, `list_saved_queries`, `get_report`, `list_reports`.

## Troubleshooting

- **Connection refused (Ollama):** Start Ollama (`ollama serve`). Docker: `LLM_BASE_URL=http://host.docker.internal:11434`. Ensure a model is pulled.
- **Report fails / timeout:** Check provider and API key; for Ollama see above. [Troubleshooting](../troubleshooting.md).
- **Slow (Ollama):** Smaller model (e.g. `mistral`), more RAM; first run per session is slower.

**See also:** [Configuration](../configuration.md), [Quick start](quickstart.md), [Troubleshooting](../troubleshooting.md)
