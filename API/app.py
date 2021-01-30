# Imports
import logging
import os
import sys
from dotenv import load_dotenv
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware


sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
load_dotenv()

# Logging on
logging.basicConfig()
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from API.db.connection import database

# Get the routers
from API.routes.pages import router as pages_router
from API.routes.users import router as users_router

# Define app and include routers and connection events.
app = FastAPI()

origins = ["*"]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

app.include_router(users_router)
app.include_router(pages_router)

# Events
@app.on_event('startup')
async def on_startup():
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    await database.disconnect()

@app.get('/')
async def welcome():
    return "THE APP IS WORKING?"

# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user