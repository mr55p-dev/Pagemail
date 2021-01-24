# Imports
from uuid import uuid4
from dotenv import load_dotenv
from datetime import datetime
import logging
import os

from sqlalchemy import Column, String, DateTime, Table
from sqlalchemy.dialects.postgresql import UUID as alchemyUUID
import sqlalchemy
import databases

from fastapi import FastAPI

# Get the environment variables
load_dotenv(verbose=True)
DATABASE_URL=os.getenv("DATABASE_URI")
SECRET = os.getenv("SECRET_KEY")

# ============================ #
#  Set up the database models  #
# ============================ #

# Database classes

# Database setup
database = databases.Database(DATABASE_URL)
metadata = sqlalchemy.MetaData()

users = Table(
    "users",
    metadata,
    Column("id", alchemyUUID(as_uuid=True), primary_key=True, default=uuid4, unique=True),
    Column("name", String),
    Column("email", String),
    Column("password_hash", String),
    Column("date_added", DateTime, server_default=sqlalchemy.sql.func.now())
    )

pages = Table(
    "pages",
    metadata,
    Column("id", alchemyUUID(as_uuid=True), primary_key=True, unique=True),
    Column("page_url", String),
    Column("date_added", DateTime, server_default=sqlalchemy.sql.func.now())
)

engine = sqlalchemy.create_engine(DATABASE_URL)
metadata.create_all(engine)

from db.models import SavePage, UserIn, UserOut, Message

# App setup
logging.basicConfig(filename="logs/application.log")
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

@app.post('/add_page', response_model=SavePage)
async def add_page_route(newPage: SavePage):
    newPage.id = uuid4()
    query = pages.insert().values(**newPage.dict())
    await database.execute(query)
    return newPage


# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user