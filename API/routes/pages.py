import logging

from datetime import datetime
from uuid import uuid4

from fastapi import APIRouter, BackgroundTasks, Form
from fastapi.exceptions import HTTPException
from fastapi.param_functions import Depends
from pydantic.error_wrappers import ValidationError

from API.db.operations import page_create, page_delete, page_verify, pages_read
from API.helpers.models import Page, UserIn
from API.helpers.utils import fetch_metadata
from API.helpers.verification import (get_partially_validated_user,
                                      get_validated_user)


router = APIRouter(
    prefix="/page",
    tags=["Pages"]
)

log = logging.getLogger(__name__)
log.setLevel(logging.DEBUG)

def generate_page(url: str = Form(""),
    current_user = Depends(get_partially_validated_user)):
    """Returns a page object from a submitted form"""
    try:
        page = Page(url=url)
        page.id = uuid4()
        page.user_id = current_user.id
        page.date_added = datetime.now()
    except ValidationError as exception:
        raise HTTPException(
        status_code=400,
        detail="the submitted URL was not valid.",
    ) from exception
    return page

@router.post('/', response_model=Page)
async def save_new_page(
    background: BackgroundTasks,
    new_page: Page = Depends(generate_page)):
    """Save a page for the user submitting it"""

    await page_create(new_page)
    background.add_task(fetch_metadata, page_id=new_page.id, url=new_page.url)

    return new_page

@router.get('/')
async def fetch_saved_pages(current_user: UserIn = Depends(get_validated_user)):
    """Fetch all pages for  user"""
    return await pages_read(current_user.id)

@router.delete('/', response_model=None)
async def delete_page_route(
        current_user: UserIn = Depends(get_validated_user),
        page_id: str = Form(None, title="page_id")):
    """Remove a page for a given user, provided they own it"""
    if (page := await page_verify(current_user, page_id)):
        # Delete metadata
        return await page_delete(page)

    raise HTTPException(
        status_code=403,
        detail="The requested page does not belong to this user"
    )

# @router.get('/fetchall')
# async def fetch_all_metadata(db: Database = Depends(get_db)):
#     scheduler.add_job(update_metadata)
#     # pages = await update_metadata(db)
#     return "Scheduled job"
