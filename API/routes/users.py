from datetime import datetime
from uuid import uuid4
from fastapi import APIRouter, Depends, HTTPException
from fastapi.security.oauth2 import OAuth2PasswordRequestForm
from api.helpers.models import User, UserOut
from api.helpers.verification import create_new_token, fetch_user, get_current_active_user, hash_password, validate_user
from api.db.connection import database, users
from asyncpg.exceptions import UniqueViolationError

router = APIRouter(
    prefix="/users",
    tags=["Users"]
)

@router.post('/add_user', response_model=UserOut)
async def add_user(new_user: User):
    new_user.id = uuid4()
    new_user.password = hash_password(new_user.password)
    new_user.date_added = datetime.utcnow()
    new_user.is_active = True

    query = users.insert().values(**new_user.dict())
    try:
        await database.execute(query=query)
    except UniqueViolationError:
        raise HTTPException(
            status_code=400,
            detail="Username already exists."
        )
    # ^Factorise this into a helper function
    return UserOut(**new_user.dict())

@router.get('/self')
async def read_self(current_user: User = Depends(get_current_active_user)):
    return current_user

@router.post('/get_token')
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    user = await fetch_user(form_data.username)
    validate_user(form_data.password, user.password)
    token = create_new_token({"sub": user.email})
    return {"access_token": token, "token_type": "bearer"}