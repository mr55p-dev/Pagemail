from pydantic.types import UUID4
from pydantic import BaseModel, EmailStr, AnyHttpUrl
from typing import Any, Optional
from datetime import datetime
from uuid import UUID, uuid4

# Database classes
class SavePage(BaseModel):
    id: Optional[UUID]
    page_url: AnyHttpUrl

class UserOut(BaseModel):
    id: Optional[UUID]
    name: Optional[str]
    email: EmailStr
    date_added: Optional[datetime]
    is_active: Optional[bool]

class User(UserOut):
    password: str

# Token classes
class TokenData(BaseModel):
    email: Optional[str] = None

# Response classes
class Message(BaseModel):
    """Send a response message"""
    message: str
    response_data: Any