# Imports
import logging
import os
import sys

from uuid import uuid4
from datetime import datetime, timedelta
from typing import Optional
from dotenv import load_dotenv

from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from starlette.types import Message


sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
load_dotenv()

# Logging on
logging.basicConfig(filename="logs/application.log")
app_log = logging.getLogger("Application Log")

# Get the database connection and models
from api.db.connection import database, users
from api.helpers.models import TokenData, User

# Verification Tools
from api.helpers.verification import get_current_active_user,\
    fetch_user, validate_user, create_new_token, verify_password

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

@app.post('/get_token')
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    # Fetch the user
    user = await fetch_user(form_data.username)
    # HASHING HERE
    # auth = verify_password(form_data.password, user.password)
    # validate_user(user, form_data.password)
    validate_user(form_data.password, user.password)
    token = create_new_token({"sub": user.email})
    return {"access_token": token, "token_type": "bearer"}

# POST: Add a user
# DELETE: Delete a user
# UPDATE: Change user info

# GET: Get all pages for a user
# UPDATE: User preferences
# DELETE: Delete a post for a user