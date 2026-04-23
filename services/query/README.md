# Query Service

Python FastAPI service for second-brain. Handles RAG-based querying over user content using Qdrant vector search and Ollama LLM.

## Setup

1. Create virtual environment:
```bash
cd services/query
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

2. Create a `.env` file in `services/query/`:
```
OLLAMA_URL=http://localhost:11434
QDRANT_URL=http://localhost:6333
COLLECTION_NAME=content_embeddings
EMBED_MODEL=nomic-embed-text
LLM_MODEL=llama3
```

3. Ensure Ollama and Qdrant are running:
```bash
ollama serve
docker run -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
```

## Run

```bash
cd services/query
source venv/bin/activate
python -m uvicorn main:app --port 8000
```

## Test

**Direct:**
```bash
curl -X POST http://localhost:8000/query -H "Content-Type: application/json" -d '{"userId":"<user_id>","query":"what videos do I have saved?"}'
```

**Via API Gateway (protected):**
```bash
curl -X POST http://localhost:8080/api/query -H "Content-Type: application/json" -H "Authorization: Bearer <token>" -d '{"query":"what videos do I have saved?"}'
```

API docs available at `http://localhost:8000/docs`
