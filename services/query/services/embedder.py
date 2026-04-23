import os
import requests

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")

def embed(text: str) -> list[float]:
    res = requests.post(
        f"{OLLAMA_URL}/api/embeddings",
        json={
            "model": "nomic-embed-text",
            "prompt": text
        },
        timeout=30
    )

    res.raise_for_status()
    return res.json()["embedding"]