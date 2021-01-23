# Imports
from typing import Optional
from uuid import UUID
from datetime import date
from dotenv import load_dotenv
import logging
import os

import sqlalchemy
import databases
from sqlalchemy.ext.declarative import DeclarativeMeta, declarative_base
from sqlalchemy import Column, String, Integer, DateTime

from fastapi import FastAPI
from fastapi_users import models
from fastapi_users.db import SQLAlchemyBaseUserTable, SQLAlchemyUserDatabase
from pydantic import BaseModel

# Get the environment variables
load_dotenv(verbose=True)
DATABASE_URL=os.getenv("DATABASE_URI")

# ============================ #
#  Set up the database models  #
# ============================ #
# User classes
class User(models.BaseUser):
    pass

class UserCreate(models.BaseUserCreate):
    pass

class UserUpdate(User, models.BaseUserUpdate):
    pass

class UserDB(User, models.BaseUserDB):
    pass

# Request/Response classes
class SavePageIn(BaseModel):
    page_url: str

class SavePage(BaseModel):
    id: UUID
    page_url: str

class Message(BaseModel):
    message: str
    page_uri: Optional[str]

# Database setup
database = databases.Database(DATABASE_URL)
base: DeclarativeMeta = declarative_base()

# Define tables
class UserTable(base, SQLAlchemyBaseUserTable):
    field_test: str
    pass

class SavedPages(base):
    __tablename__ = "pages"
    id = Column(String, primary_key=True)
    url = Column(String)
    date_added = Column(DateTime)

users = UserTable.__table__
user_db = SQLAlchemyUserDatabase(UserDB, database, User)

# Generating the tables
# Add in a migration system (alembic)
engine = sqlalchemy.create_engine(DATABASE_URL)
base.metadata.create_all(engine)

# ==================== #
# Authentication setup #
# ==================== #
from fastapi_users.authentication import JWTAuthentication

SECRET = os.getenv("SECRET_KEY")
jwt_auth = JWTAuthentication(SECRET, lifetime_seconds=3600)

# =========== #
# Users setup #
# =========== #
from fastapi_users import FastAPIUsers

fapi_users = FastAPIUsers(
    user_db,
    [jwt_auth],
    User,
    UserCreate,
    UserUpdate,
    UserDB
)

# App setup
logging.basicConfig()
app_log = logging.getLogger("Application Log")

app = FastAPI()
app.include_router(fapi_users.get_auth_router(jwt_auth), prefix="/auth/jwt", tags=["auth"])
app.include_router(fapi_users.get_register_router(), prefix="/auth", tags=["auth"])

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

# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user