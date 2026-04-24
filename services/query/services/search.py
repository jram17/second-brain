import os
from qdrant_client import QdrantClient, models
import pybreaker
cb = pybreaker.CircuitBreaker(fail_max=3, reset_timeout=30)


QDRANT_URL = os.getenv("QDRANT_URL", "http://localhost:6333")
COLLECTION = os.getenv("COLLECTION_NAME", "content_embeddings")
client = QdrantClient(url=QDRANT_URL)

@cb
def search(vector: list[float], user_id: str, limit: int = 5):
    results = client.query_points(
        collection_name=COLLECTION,
        query=vector,
        limit=limit,
        query_filter=models.Filter(
            must=[
                models.FieldCondition(
                    key="userId",
                    match=models.MatchValue(value=user_id)
                )
            ]
        )
    )

    return [
        {
            "contentId": r.payload.get("contentId", ""),
            "text": r.payload.get("text", ""),
            "score": r.score
        }
        for r in results.points
    ]
