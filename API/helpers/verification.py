import os
from typing import Optional
from datetime import datetime, timedelta
from jose.exceptions import ExpiredSignatureError
from passlib.context import CryptContext
from jose import JWTError, jwt
from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer

from API.helpers.models import UserIn, TokenData
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
async def fetch_user(email: str) -> Optional[UserIn]:
    query = users.select().where(users.c.email == email)
    user = await database.fetch_one(query=query)
    if not user:
        raise HTTPException(
            status_code=400,
            detail="Incorrect user name."
            )
    return UserIn(**user)

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
        detail="the login token could not be validated.",
        headers={"WWW-Authenticate": "Bearer"}
    )

    expired_token_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="the login token has expired.",
        headers={"WWW-Authenticate": "Bearer"}
    )

    try:
        payload = jwt.decode(token, SECRET, ALGORITHM)
        username: str = payload.get("sub")
        scope: str = payload.get("scope")
        if (username or scope) is None:
            raise cred_exception
        token = TokenData(email=username, scope=scope)
    except JWTError:
        raise cred_exception
    except ExpiredSignatureError:
        raise expired_token_exception

    return token

def create_new_token(data: dict, expires_delta: Optional[timedelta] = None, page_only: bool = False):
    to_encode = data.copy()
    # if not expires_delta:
    #     expires = datetime.utcnow() + timedelta(minutes=TOKEN_EXPIRATION_TIME)
    # elif expires_delta > 0:
    #     expires = datetime.utcnow() + expires_delta
    # else:
    #     expires = datetime.max
    if page_only:
        expires = datetime.max
    elif expires_delta:
        expires = datetime.utcnow() + timedelta(minutes=expires_delta)
    else:
        expires = datetime.utcnow() + timedelta(minutes=TOKEN_EXPIRATION_TIME)

    to_encode.update({"exp": expires})
    to_encode.update({"scope": "userauth:none" if page_only else "userauth:full"})
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
    try:
        if token.scope == "userauth:full":
            auth = 1
        elif token.scope == "userauth:none":
            auth = 0
        else:
            raise AttributeError
    except AttributeError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="an invalid token was submitted, please refresh your login."
        )

    return (UserIn(**user.dict()), auth)

async def get_current_active_user(userauth = Depends(get_current_user)):
    user, auth = userauth
    if not user.is_active:
        raise HTTPException(
            status_code=400,
            detail="this user is not marked active currently, please contact the site"
        )
    return (user, auth)

async def get_validated_user(userauth = Depends(get_current_active_user)):
    user, auth = userauth
    if auth != 1:
        raise HTTPException(
            status_code=401,
            detail="the login credentials are not sufficent, please sign in from a browser."
        )
    else:
        return user

async def get_partially_validated_user(userauth = Depends(get_current_active_user)):
    user, auth = userauth
    if auth >= 0:
        return user
    else:
        raise HTTPException(
            status=403,
            detail="the permissions of your login could not be verified."
        )
