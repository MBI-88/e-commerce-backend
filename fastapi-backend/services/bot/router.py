from middlewares.jwt import authenticate
from fastapi import APIRouter, HTTPException, Response, status, Depends
from models.bot import Bot


router = APIRouter()

# Routes

@router.get("/home")
async def getBot() -> object:
    return Response(content=str({"message":"Hola mundo"}),status_code=status.HTTP_200_OK)


@router.post("/bot")
async def messageBot(user: dict = Depends(authenticate)) -> object:
    return Response(content=str({"message":user}),status_code=status.HTTP_200_OK)