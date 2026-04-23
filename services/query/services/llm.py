import os
import requests

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")

def generate_answer(query: str, contexts: list[str]) -> str:
    context_text = "\n\n".join(contexts)

    prompt = f"""
You are a helpful assistant.

Use the following context to answer the question.

Context:
{context_text}

Question:
{query}

Answer:
"""

    res = requests.post(
        f"{OLLAMA_URL}/api/generate",
        json={
            "model": "llama3",
            "prompt": prompt,
            "stream": False
        },
        timeout=60
    )

    res.raise_for_status()
    return res.json()["response"]