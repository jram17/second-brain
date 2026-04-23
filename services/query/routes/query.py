from fastapi import APIRouter,HTTPException
from pydantic import BaseModel

from services.embedder import embed
from services.search import search
from services.llm import generate_answer


router=APIRouter()

class QueryRequest(BaseModel):
    userId: str
    query: str

class QuerResponse(BaseModel):
    answer: str
    sources: list

@router.post("/query",response_model=QuerResponse)
def query_endpoint(req: QueryRequest):
    try:
        vector = embed(req.query)
        results = search(vector, req.userId)
        contexts = [r["text"] for r in results if r["text"]]
        answer= generate_answer(req.query,contexts)
        return {
            "answer": answer,
            "sources":results
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))