# Imports
import logging
import os
import sys
from uuid import uuid4

from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm

from datetime import datetime, timedelta

import sqlalchemy
from sqlalchemy.sql.expression import desc
from typing import Optional
from dotenv import load_dotenv

sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
load_dotenv()
# Logging on
logging.basicConfig(filename="logs/application.log")
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from api.db.connection import database, pages, users
from api.db.models import SavePage, TokenData, User, Message

# Get the routers
from api.routes.pages import router as pages_router
from api.routes.users import router as users_router

# Define app and include routers and connection events.
app = FastAPI()
app.include_router(users_router)
app.include_router(pages_router)

# Security begins

from passlib.context import CryptContext
from jose import JWTError, jwt
SECRET = os.getenv("SECRET_KEY")
ALGORITHM = os.getenv("ALGORITHM")
TOKEN_EXPIRATION_TIME = int(os.getenv("TOKEN_DURATION"))

password_context = CryptContext(schemes=["bcrypt"], deprecated="auto")
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

def hash_password(password: str):
    return password_context.hash(password)

def verify_password(plain: str, hash: str):
    return password_context.verify(plain, hash)

def create_new_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    if expires_delta:
        expires = datetime.utcnow() + expires_delta
    else:
        expires = datetime.utcnow() + timedelta(minutes=TOKEN_EXPIRATION_TIME)
    to_encode.update({"exp": expires})
    encoded = jwt.encode(to_encode, SECRET, ALGORITHM)
    return encoded


async def fetch_user(email: str) -> Optional[User]:
    query = users.select().where(users.c.email == email)
    return await database.fetch_one(query=query)


async def decode_token(token: str):
    # This will actually do something important one day.
    cred_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate token credientials.",
        headers={"WWW-Authenticate": "Bearer"}
    )
    try:
        payload = jwt.decode(token, SECRET, ALGORITHM)
        username: str = payload.get("sub")
        if username is None:
            raise cred_exception
        token = TokenData(email=username)
    except JWTError:
        raise cred_exception

    return token

async def get_current_user(token: str = Depends(oauth2_scheme)):
    token = await decode_token(token)
    user = await fetch_user(token.email)
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
    user = await fetch_user(form_data.username)
    if not user:
        raise HTTPException(
            status_code=400,
            detail="Incorrect user name."
            )
    user = User(**user)
    # Move ^ this logic into the fetch_user function.
    # HASHING HERE
    # auth = verify_password(form_data.password, user.password)
    if False:
        raise HTTPException(
            status_code=400,
            detail="Incorrect password."
            )
    # Create a validate user function for ^ this logic.
    token = create_new_token({"sub": user.email})
    return {"access_token": token, "token_type": "bearer"}

@app.get('/user/self')
async def read_users_self(current_user: User = Depends(get_current_active_user)):
    return current_user

# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user