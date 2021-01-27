# Imports
import logging
import os
import sys
from uuid import uuid4

from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm

from datetime import datetime

import sqlalchemy
from sqlalchemy.sql.expression import desc
from typing import Optional

sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))

# Logging on
logging.basicConfig(filename="logs/application.log")
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from api.db.connection import database, pages, users
from api.db.models import SavePage, User, Message

# Get the routers
from api.routes.v1.users import router as users_router
from api.routes.v1.pages import router as pages_router

# Define app and include routers and connection events.
app = FastAPI()
app.include_router(users_router)
app.include_router(pages_router)

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

async def fetch_user(email: str) -> Optional[User]:
    query = users.select().where(users.c.email == email)
    return await database.fetch_one(query=query)


async def decode_token(token: str):
    # This will actually do something important one day.
    return await fetch_user(token)

async def get_current_user(token: str = Depends(oauth2_scheme)):
    user = await decode_token(token)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid crediantials",
            headers={"WWW-Authenticate": "Bearer"}
        )

    return User(**user)

async def get_current_active_user(user: User = Depends(get_current_user)):
    if not user.is_active:
        raise HTTPException(
            status_code=400,
            detail="User is not active"
        )
    return user

@app.on_event('startup')
async def on_startup():
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    await database.disconnect()

@app.post('/token')
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    # Fetch the user
    query = users.select().where(users.c.email == form_data.username)
    user = await fetch_user(form_data.username)

    if not user:
        raise HTTPException(
            status_code=400,
            detail="Incorrect user name."
            )
    user = User(**user)
    # HASHING HERE
    if user.password != form_data.password:
        raise HTTPException(
            status_code=400,
            detail="Incorrect password."
            )

    return {"access_token": user.email, "token_type": "bearer"}

@app.get('/user/self')
async def read_users_self(current_user: User = Depends(get_current_active_user)):
    return current_user

# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user