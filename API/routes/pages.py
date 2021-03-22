from typing import List
from uuid import uuid4
from databases.core import Database
from fastapi import APIRouter, BackgroundTasks, Form
from fastapi.exceptions import HTTPException
from fastapi.openapi.models import HTTPBase
from fastapi.param_functions import Depends
from API.db.connection import get_db, pages, users, page_metadata
from API.helpers.models import Page, PageFilled, PageMetadata, PageOut, UserIn, UserOut
from API.helpers.scheduling import scheduler
from API.helpers.utils import fetch_page, fetch_pages, set_page_metadata, unwrap_submitted_page, fetch_metadata, update_metadata, verify_ownership
from API.helpers.verification import get_current_active_user
from sqlalchemy.sql import select
from datetime import datetime

router = APIRouter(
    prefix="/page",
    tags=["Pages"]
)

def page_saved(response):
    # if not response:
    #     print("Something went wrong")
    # else:
    #     print(f"Apparently all good: {response}")
    pass

@router.post('/save', response_model=Page)
async def save_new_page(
    background: BackgroundTasks,
    new_page: Page = Depends(unwrap_submitted_page),
    current_user: UserIn = Depends(get_current_active_user),
    database: Database = Depends(get_db)
    ):
    
    new_page.id = uuid4()
    new_page.user_id = current_user.id
    new_page.date_added = datetime.now()

    query = pages.insert().values(**new_page.dict())
    await database.execute(query=query)
    background.add_task(fetch_metadata, id=new_page.id, url=new_page.url, db=database, callback=page_saved)

    return new_page

@router.get('/mypages')
async def fetch_saved_pages(
    current_user: UserIn = Depends(get_current_active_user),
    database: Database = Depends(get_db)):

    query = users.join(pages)
    query = select([pages]).select_from(query).where(users.c.id == current_user.id)
    results = await database.fetch_all(query=query)
    user_pages = [PageFilled(**i) for i in results]
    page_ids = [i.id for i in user_pages]

    query = pages.join(page_metadata)
    query = select([page_metadata]).select_from(query).where(page_metadata.c.id.in_(page_ids))
    results_meta = await database.fetch_all(query=query)
    pages_meta = [PageMetadata(**i) for i in results_meta]

    return await set_page_metadata(user_pages, pages_meta)

# NEed to change this to a way that we can get the user, the page and also if the user owns the page...
# Maybe just a veryify function...
@router.delete('/delete', response_model=None)
async def delete_page(
    current_user: UserIn = Depends(get_current_active_user),
    id: str = Form(None, title="id"),
    db: Database = Depends(get_db)):

    if await verify_ownership(current_user.id, id):
        # Delete metadata
        query = page_metadata.delete().where(page_metadata.c.id == id)
        response = await db.execute(query)
        # Delete page
        query = pages.delete().where(pages.c.id == id)
        response = await db.execute(query)
        return 200
    else:
        raise HTTPException(
            status_code=403,
            detail="the requested page does not belong to this user"
        )

@router.get('/fetchall')
async def fetch_all_metadata(db: Database = Depends(get_db)):
    scheduler.add_job(update_metadata)
    # pages = await update_metadata(db)
    return "Scheduled job"
