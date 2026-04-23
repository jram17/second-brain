from fastapi import FastAPI
from dotenv import load_dotenv
from routes.query import router as query_router
load_dotenv()
app=FastAPI(title="Query service")
app.include_router(query_router)