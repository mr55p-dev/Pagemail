import os
from datetime import datetime, timedelta
from typing import Optional

from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from jose.exceptions import ExpiredSignatureError
from passlib.context import CryptContext

from API.db.operations import user_read
from API.helpers.models import TokenData, UserIn

SECRET = os.getenv("SECRET_KEY")
ALGORITHM = os.getenv("ALGORITHM")
TOKEN_EXPIRATION_TIME = int(os.getenv("TOKEN_DURATION"))

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")
password_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

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

# Password functions
def hash_password(password: str):
    """Returns a password hash"""
    return password_context.hash(password)

def verify_password(plain: str, pass_hash: str):
    """Verifies a password against a saved hash"""
    return password_context.verify(plain, pass_hash)

# User functions
async def fetch_user(email: str) -> Optional[UserIn]:
    """Fetches a user from the table via email"""
    return await user_read(email)

def validate_user(submitted_password: str, hashed_password: str) -> bool:
    """Catches validation errors and returns a boolean for password verification"""
    if not verify_password(submitted_password, hashed_password):
        raise HTTPException(
            status_code=400,
            detail="Incorrect password."
            )
    return True

# Token functions
async def decode_token(token: str):
    """Decode a submitted JWT token"""
    # This will actually do something important one day.
    try:
        payload = jwt.decode(token, SECRET, ALGORITHM)
        username: str = payload.get("sub")
        scope: str = payload.get("scope")
    except ExpiredSignatureError as exception:
        raise expired_token_exception from exception
    except JWTError as exception:
        raise cred_exception from exception
    else:
        if (username or scope) is None:
            raise cred_exception
        token = TokenData(email=username, scope=scope)

    return token

def create_new_token(
        data: dict,
        expires_delta: Optional[timedelta] = None,
        page_only: bool = False):
    """Creates a token with the given permission and expiry"""
    to_encode = data.copy()
    if page_only:
        expires = datetime.max
    elif expires_delta:
        expires = datetime.utcnow() + timedelta(minutes=expires_delta)
    else:
        expires = datetime.utcnow() + timedelta(minutes=TOKEN_EXPIRATION_TIME)

    to_encode.update({"exp": expires})
    to_encode.update({"scope": "userauth:none" if page_only else "userauth:full"})
    return jwt.encode(to_encode, SECRET, ALGORITHM)

# User and token functions
async def get_current_user(token: str = Depends(oauth2_scheme)) -> UserIn:
    """Converts a token into a user object"""
    token = await decode_token(token)
    user = await user_read(user_email=token.email)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid crediantials",
            headers={"WWW-Authenticate": "Bearer"}
        )
    if token.scope == "userauth:full":
        auth = 1
    elif token.scope == "userauth:none":
        auth = 0
    else:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="An invalid token was submitted, please refresh your login."
        )

    return (UserIn(**user.dict()), auth)

async def get_current_active_user(userauth = Depends(get_current_user)):
    """Catches inactive users and throws them as an error"""
    user, auth = userauth
    if not user.is_active:
        raise HTTPException(
            status_code=400,
            detail="this user is not marked active currently, please contact the site"
        )
    return (user, auth)

async def get_validated_user(userauth = Depends(get_current_active_user)):
    """Catches if a token has sufficient permission"""
    user, auth = userauth
    if auth != 1:
        raise HTTPException(
            status_code=401,
            detail="the login credentials are not sufficent, please sign in from a browser."
        )
    else:
        return user

async def get_partially_validated_user(userauth = Depends(get_current_active_user)):
    """Ensures tokens have any authentication at all"""
    user, auth = userauth
    if auth >= 0:
        return user
    else:
        raise HTTPException(
            status_code=403,
            detail="the permissions of your login could not be verified."
        )
