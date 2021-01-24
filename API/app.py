# Imports
from uuid import uuid4
from dotenv import load_dotenv
import logging
import os
import sys

import sqlalchemy
import databases
from fastapi import FastAPI

sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))

# Logging on
logging.basicConfig(filename="logs/application.log")
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from api.db.connection import database, pages, users
from api.db.models import SavePage, UserIn, UserOut, Message

# Get the routers
from api.routes.v1.users import router as users_router
from api.routes.v1.pages import router as pages_router

# Define app and include routers and connection events.
app = FastAPI()
app.include_router(users_router)
app.include_router(pages_router)

@app.on_event('startup')
async def on_startup():
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    await database.disconnect()


# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user