# Imports
import logging
import os
import sys
from dotenv import load_dotenv
from fastapi import FastAPI


sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
load_dotenv()

# Logging on
logging.basicConfig()
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from api.db.connection import database

# Get the routers
from api.routes.pages import router as pages_router
from api.routes.users import router as users_router

# Define app and include routers and connection events.
app = FastAPI()
app.include_router(users_router)
app.include_router(pages_router)

# Events
@app.on_event('startup')
async def on_startup():
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    await database.disconnect()


# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user