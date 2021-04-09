from datetime import datetime
from typing import Any, List, Optional, Union
from uuid import UUID

from pydantic import AnyHttpUrl, BaseModel, EmailStr

# Database classes

class PageOut(BaseModel):
    """Pages which are sanitised to not contain user_id information"""
    id: Optional[UUID]
    url: AnyHttpUrl
    date_added: Optional[datetime]
    # title: Optional[str]
    # description: Optional[str]

class Page(PageOut):
    """Pages for internal use with user_id information"""
    user_id: Optional[UUID]

class PageMetadata(BaseModel):
    """Metadata object not associated with a page directly"""
    id: Optional[UUID]
    title: str
    description: str

class PageFilled(Page):
    """Page combining all the peoperties of metadata and page"""
    title: Optional[str]
    description: Optional[str]

class UserOut(BaseModel):
    """User for external use, with no password information"""
    id: UUID
    name: str
    email: EmailStr
    date_added: Optional[datetime]
    is_active: Optional[bool]
    token: Optional[str]

class UserIn(BaseModel):
    """User for internal use which carries active and password information"""
    id: Optional[UUID]
    name: Optional[str]
    email: EmailStr
    password: str
    date_added: Optional[datetime]
    is_active: Optional[bool]


# Token classes
class TokenData(BaseModel):
    """Class which allows mainpulating token information"""
    email: str
    scope: str

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
