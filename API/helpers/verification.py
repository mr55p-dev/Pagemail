import os
from typing import Optional
from datetime import datetime, timedelta
from passlib.context import CryptContext
from jose import JWTError, jwt
from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer

from API.helpers.models import User, TokenData
from API.db.connection import users, database
from fastapi import HTTPException, status

SECRET = os.getenv("SECRET_KEY")
ALGORITHM = os.getenv("ALGORITHM")
TOKEN_EXPIRATION_TIME = int(os.getenv("TOKEN_DURATION"))

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")
password_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

# Password functions
def hash_password(password: str):
    return password_context.hash(password)

def verify_password(plain: str, hash: str):
    return password_context.verify(plain, hash)

# User functions
async def fetch_user(email: str) -> Optional[User]:
    query = users.select().where(users.c.email == email)
    user = await database.fetch_one(query=query)
    if not user:
        raise HTTPException(
            status_code=400,
            detail="Incorrect user name."
            )
    return User(**user)

def validate_user(submitted_password: str, hashed_password: str) -> bool:
    if not verify_password(submitted_password, hashed_password):
        raise HTTPException(
            status_code=400,
            detail="Incorrect password."
            )
    return True

# Token functions
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

def create_new_token(data: dict, expires_delta: Optional[timedelta] = None):
    to_encode = data.copy()
    if expires_delta:
        expires = datetime.utcnow() + expires_delta
    else:
        expires = datetime.utcnow() + timedelta(minutes=TOKEN_EXPIRATION_TIME)
    to_encode.update({"exp": expires})
    encoded = jwt.encode(to_encode, SECRET, ALGORITHM)
    return encoded

# User and token functions
async def get_current_user(token: str = Depends(oauth2_scheme)):
    token = await decode_token(token)
    user = await fetch_user(token.email)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid crediantials",
            headers={"WWW-Authenticate": "Bearer"}
        )

    return User(**user.dict())

async def get_current_active_user(user: User = Depends(get_current_user)):
    if not user.is_active:
        raise HTTPException(
            status_code=400,
            detail="User is not active"
        )
    return user
