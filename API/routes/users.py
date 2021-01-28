from uuid import uuid4
from fastapi import APIRouter, Depends
from api.helpers.models import User, UserOut
from api.helpers.verification import get_current_active_user, hash_password
from api.db.connection import database, users

router = APIRouter(
    prefix="/users",
    tags=["Users"]
)

@router.post('/add_user', response_model=UserOut)
async def add_user(new_user: User):
    new_user.id = uuid4()
    new_user.password = hash_password(new_user.password)

    query = users.insert().values(**new_user.dict())
    await database.execute(query=query)
    # ^Factorise this into a helper function
    return UserOut(**new_user.dict())

@router.get('/self')
async def read_self(current_user: User = Depends(get_current_active_user)):
    return current_user