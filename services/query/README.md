# Query Service

Python FastAPI service. RAG-based querying over user content using Qdrant vector search + Ollama LLM. Circuit breakers on all external calls.

## Setup

```bash
cd services/query
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

Create `.env`:
```
OLLAMA_URL=http://localhost:11434
QDRANT_URL=http://localhost:6333
COLLECTION_NAME=content_embeddings
EMBED_MODEL=nomic-embed-text
LLM_MODEL=llama3
```

Requires Qdrant and Ollama running.

## Run

```bash
source venv/bin/activate
python -m uvicorn main:app --port 8000
```

## How it works

```
Query → embed query (Ollama) → vector search (Qdrant, filtered by userId) → build prompt with context → LLM generates answer (Ollama) → return answer + sources
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /query | Query content with RAG |

Auto-generated docs: `http://localhost:8000/docs`

## Test

```bash
# Direct
curl -X POST http://localhost:8000/query \
  -H "Content-Type: application/json" \
  -d '{"userId":"<user_id>","query":"what videos do I have saved?"}'

# Via Gateway (protected)
curl -X POST http://localhost:8080/api/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"query":"what videos do I have saved?"}'
```
