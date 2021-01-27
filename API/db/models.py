from pydantic.types import UUID4
from pydantic import BaseModel, EmailStr, AnyHttpUrl
from typing import Any, Optional
from datetime import datetime
from uuid import UUID, uuid4

# Database classes
class SavePage(BaseModel):
    id: Optional[UUID]
    page_url: AnyHttpUrl

class User(BaseModel):
    id: Optional[UUID]
    name: str
    email: EmailStr
    password: str
    date_added: datetime
    is_active: bool


# Response classes
class Message(BaseModel):
    """Send a response message"""
    message: str
    response_data: Any