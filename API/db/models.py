from pydantic.types import UUID4
from pydantic import BaseModel, EmailStr, AnyHttpUrl
from typing import Any, Optional
from datetime import datetime
from uuid import UUID, uuid4

# Database classes
class SavePage(BaseModel):
    id: Optional[UUID]
    page_url: AnyHttpUrl

class UserIn(BaseModel):
    name: str
    email: EmailStr
    password_hash: str
    date_added: datetime

class UserOut(BaseModel):
    id: UUID
    name: str
    email: EmailStr
    date_added: datetime

# Response classes
class Message(BaseModel):
    """Send a response message"""
    message: str
    response_data: Any