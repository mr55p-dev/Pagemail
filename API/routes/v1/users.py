from fastapi import APIRouter
from api.db.connection import database, users

router = APIRouter(
    prefix="/v1/users",
    tags=["Users"]
)

@router.post('/add_user')
def add_user():
    return {"message": "cool"}