# Imports
import logging
import os
import sys
from dotenv import load_dotenv
from datetime import timedelta
from fastapi.middleware.cors import CORSMiddleware
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

sys.path.append(os.path.dirname(os.path.dirname(os.path.realpath(__file__))))
load_dotenv()

# Logging on
# log_format_str = r"%(name)s: %(lineno)d: %(msg)s"
logging.basicConfig()

# Change in production
# logging.getLogger('uvicorn').setLevel(logging.ERROR)
# logging.getLogger('apscheduler').setLevel(logging.INFO)
logging.getLogger('').setLevel(logging.DEBUG)

# Get the database connection and models
from API.db.connection import database

# Get the task scheduler and start it on app launch
from async_scheduler.job import BaseJob
from API.helpers.scheduling import scheduler, sch_news, sch_meta
from API.helpers.email_tools import newsletter

# Get the jobs which are to be scheduled on startup
from API.helpers.utils import update_metadata

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
    scheduler.start()
    METADATA_UPDATE_INTERVAL = int(os.getenv("METADATA_UPDATE_INTERVAL")) or 60
    # scheduler.add_job(update_metadata, 'interval', minutes=METADATA_UPDATE_INTERVAL, id="1", jobstore="local")

    sch_meta.function = update_metadata
    sch_news.function = newsletter

    interval = timedelta(minutes=METADATA_UPDATE_INTERVAL)
    update_job = BaseJob(interval, name="update metadata", id=1)
    sch_meta.add(update_job)

    sch_meta.start()
    # my_scheduler.start()

    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    # scheduler.shutdown()
    sch_news.stop()
    sch_meta.stop()
    await database.disconnect()

@app.get('/')
async def welcome():
    return "Hello World."
