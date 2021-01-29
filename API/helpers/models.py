from pydantic.types import UUID4
from pydantic import BaseModel, EmailStr, AnyHttpUrl
from typing import Any, Optional
from datetime import datetime
from uuid import UUID, uuid4

# Database classes
class Page(BaseModel):
    id: Optional[UUID]
    page_url: AnyHttpUrl

class UserOut(BaseModel):
    name: str
    email: EmailStr
    date_added: Optional[datetime]


class User(UserOut):
    id: Optional[UUID]
    password: str
    is_active: Optional[bool]


# Token classes
class TokenData(BaseModel):
    email: Optional[str] = None

# Response classes
class Message(BaseModel):
    """Send a response message"""
    message: str
    response_data: Any