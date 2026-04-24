import os
import requests
import pybreaker

cb=pybreaker.CircuitBreaker(fail_max=3, reset_timeout=30)       


OLLAMA_URL = os.getenv("OLLAMA_URL", "http://localhost:11434")

@cb
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