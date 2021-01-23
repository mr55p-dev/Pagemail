# Imports
from typing import Optional
from uuid import UUID, uuid4
from datetime import datetime
from dotenv import load_dotenv
import logging
import os
from fastapi.params import Form
from pydantic.networks import AnyHttpUrl, EmailStr

import sqlalchemy
import databases
from sqlalchemy.ext.declarative import DeclarativeMeta, declarative_base
from sqlalchemy import Column, String, DateTime

from fastapi import FastAPI
from pydantic import BaseModel

# Get the environment variables
load_dotenv(verbose=True)
DATABASE_URL=os.getenv("DATABASE_URI")
SECRET = os.getenv("SECRET_KEY")

# ============================ #
#  Set up the database models  #
# ============================ #

# Database classes
class SavePage(BaseModel):
    id: UUID
    page_url: AnyHttpUrl

class UserIn(BaseModel):
    id: UUID = uuid4()
    name: str
    email: EmailStr
    password_hash: str
    date_added: datetime = datetime.now()

class UserOut(BaseModel):
    id: UUID
    name: str
    email: EmailStr
    date_added: datetime

# Response classes

class Message(BaseModel):
    """Send a response message"""
    message: str
    page_uri: Optional[str]

# Database setup
database = databases.Database(DATABASE_URL)
base: DeclarativeMeta = declarative_base()

# User tables
class UserTable(base):
    __tablename__ = "users"
    id = Column(String, primary_key=True)
    name = Column(String)
    email = Column(String)
    password_hash = Column(String)
    date_added = Column(DateTime)
    # Need a link to owned pages.

class SavedPages(base):
    __tablename__ = "pages"
    id = Column(String, primary_key=True)
    url = Column(String)
    date_added = Column(DateTime)

# Generating the tables
# Add in a migration system (alembic)
engine = sqlalchemy.create_engine(DATABASE_URL)
base.metadata.create_all(engine)


# App setup
logging.basicConfig()
app_log = logging.getLogger("Application Log")

app = FastAPI()

@app.on_event('startup')
async def on_startup():
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    await database.disconnect()


# Routes
@app.get('/', response_model=Message)
async def root_route():
    return Message(message="Hello World")

@app.post('/add_page', response_model=Message)
async def add_page_route(page_uri: str):
    app_log.info(f"Received page URI {page_uri}")
    return Message(
        message="Recieved the page",
        page_uri=page_uri
        )

@app.post('/user/register', response_model=UserOut)
async def register_route(newUser: UserIn):
    return newUser

# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user