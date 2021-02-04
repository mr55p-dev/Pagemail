from typing import List
from uuid import uuid4
from databases.core import Database
from fastapi import APIRouter
from fastapi.param_functions import Depends
from API.db.connection import database, get_db, pages, users
from API.helpers.models import Page, PageOut, UserIn, UserOut
from API.helpers.utils import unwrap_submitted_page
from API.helpers.verification import get_current_active_user
from sqlalchemy.sql import select
from datetime import datetime

router = APIRouter(
    prefix="/page",
    tags=["Pages"]
)

@router.post('/save', response_model=Page)
async def save_new_page(
    new_page: Page = Depends(unwrap_submitted_page),
    current_user: UserIn = Depends(get_current_active_user),
    database: Database = Depends(get_db)):
    
    new_page.id = uuid4()
    new_page.user_id = current_user.id
    new_page.date_added = datetime.now()

    query = pages.insert().values(**new_page.dict())
    await database.execute(query=query)
    return new_page

@router.get('/mypages', response_model=List[PageOut])
async def fetch_saved_pages(
    current_user: UserIn = Depends(get_current_active_user),
    database: Database = Depends(get_db)):

    query = users.join(pages)
    query = select([pages]).select_from(query).where(users.c.id == current_user.id)
    results = await database.fetch_all(query=query)
    return [PageOut(**i) for i in results]
