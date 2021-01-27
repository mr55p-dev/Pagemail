from uuid import uuid4

from pydantic.errors import UrlUserInfoError
from api.db.models import User
from fastapi import APIRouter
from api.db.connection import database, users

router = APIRouter(
    prefix="/v1/users",
    tags=["Users"]
)

@router.post('/add_user')
async def add_user(new_user: User):
    new_user.id = uuid4()
    query = users.insert().values(**new_user.dict())
    await database.execute(query=query)
    return new_user