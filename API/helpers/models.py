from pydantic.types import UUID4
from pydantic import BaseModel, EmailStr, AnyHttpUrl
from typing import Any, List, Optional, Union
from datetime import datetime
from uuid import UUID, uuid4

# Database classes

class PageOut(BaseModel):
    id: Optional[UUID]
    url: AnyHttpUrl
    date_added: Optional[datetime]
    # title: Optional[str]
    # description: Optional[str]

class Page(PageOut):
    user_id: Optional[UUID]

class PageMetadata(BaseModel):
    id: Optional[UUID]
    title: str
    description: str

class PageFilled(PageOut):
    title: Optional[str]
    description: Optional[str]

class UserOut(BaseModel):
    id: UUID
    name: str
    email: EmailStr
    date_added: Optional[datetime]

class UserIn(BaseModel):
    id: Optional[UUID]
    name: Optional[str]
    email: EmailStr
    password: str
    date_added: Optional[datetime]
    is_active: Optional[bool]


# Token classes
class TokenData(BaseModel):
    email: Optional[str] = None

# Email classes
class BaseEmail(BaseModel):
    """
    sender: email
    recipients: email | List[email]
    subject: str
    content: str
    """

    sender: EmailStr = "noreply@pagemail.tech"
    subject: str
    content: str
    recipients: Union[EmailStr, List[EmailStr]]

# Response classes
class Message(BaseModel):
    """Send a response message"""
    message: str
    response_data: Any