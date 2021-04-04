# Imports
import logging
import os
import sys
from dotenv import load_dotenv
from fastapi.middleware.cors import CORSMiddleware
from fastapi import FastAPI

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
from API.helpers.scheduling import scheduler, my_scheduler
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
    METADATA_UPDATE_INTERVAL = int(os.getenv("METADATA_UPDATE_INTERVAL")) if os.getenv("METADATA_UPDATE_INTERVAL") else 60
    scheduler.add_job(update_metadata, 'interval', minutes=METADATA_UPDATE_INTERVAL, id="1", jobstore="local")

    my_scheduler.function = newsletter
    my_scheduler.start()
    await database.connect()

@app.on_event('shutdown')
async def on_shutdown():
    scheduler.shutdown()
    await database.disconnect()

@app.get('/')
async def welcome():
    return "Hello World."

"""
helpers/utils.py -> helpers/server_actions.py:
            update_metadata()
            fetch_metadata()
            fetch_page() MOVE into REQUESTS
            verify_page_ownership(user_id, page_id)
            unwrap_page(url: str = Form(""))

helpers/verification.py ->
            decode_new_user_form,
            decode_user_form


db/requests.py ->
            create_user(new_user)
            read_user(user_id)
            update_user(user_id, fields)
            delete_user(user_id)

            create_page(new_page)
            read_page(user_id?page_id, all: bool) -> Union[Page, List[Page]]
            update_page(user_id, fields)
            delete_page(user_id, all: bool)

            create_metadata(page_id)
            read_metadata(page_id)
            update_metadata(page_id, fields)
            delete_metadata(page_id)

"""