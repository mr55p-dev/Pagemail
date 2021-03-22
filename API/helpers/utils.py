from uuid import UUID

from sqlalchemy.sql import select
from API.helpers.models import Page
from API.helpers.models import PageMetadata
from API.db.connection import get_db, page_metadata, pages
from fastapi import Form, HTTPException, Depends
from pydantic.error_wrappers import ValidationError
from bs4 import BeautifulSoup as BS
from databases.core import Database

import requests


async def update_metadata():
    db = next(get_db())
    query = pages.join(page_metadata, full=True)
    query = select([pages.c.id, pages.c.url]).select_from(query).where(page_metadata.c.id == None)
    unannotated_pages = await db.fetch_all(query)
    for page in unannotated_pages:
        await fetch_metadata(id=page["id"], url=page["url"], db=db)
    return unannotated_pages

# This should be redundent
async def fetch_pages(id, db = Depends(get_db)):
    query = pages.select().where(id == pages.c.id)
    result = await db.fetch_all(query)
    if not result:
        raise HTTPException(
            status_code=404,
            detail="the page could not be found."
        )
    return result

async def fetch_page(id: int = Form(...), db = Depends(get_db)):
    query = pages.select().where(id == pages.c.id)
    result = await db.fetch_one(query)
    if not result:
        raise HTTPException(
            status_code=404,
            detail="the page could not be found."
        )
    return Page(**result)

# Can merge all these page db requests into one or two
async def verify_ownership(user_id, page_id, db):
        query = select([pages.c.user_id]).select_from(pages).where(page_id == pages.c.id)
        result = await db.fetch_one(query)
        if result == user_id:
            return True
        else:
            return False

async def unwrap_submitted_page(url: str = Form("")):
    try:
        page = Page(url=url)
    except ValidationError as e:
        raise HTTPException(
        status_code=400,
        detail="the submitted URL was not valid.",
    )
    return page

async def fetch_metadata(id: UUID, url: str, db: Database, callback: callable = print):
    try:
        resp = requests.get(url)
    except:
        title = ""
        desc = ""
    else:
        # Validate response
        query = pages.select().where(id == pages.c.id)
        page = await db.fetch_one(query=query)
        if not page:
            callback(ValidationError('The requested page does not exist.'))
            return None

        # Get page title
        soup = BS(resp.text, features="html.parser")

        title = soup.title.text
        desc = soup.get_text()[len(title):500]

    meta = PageMetadata(
        id=id,
        title=title,
        description=desc
    )
    query = page_metadata.insert().values(**meta.dict())
    result = await db.execute(query)
    # callback(result)

async def set_page_metadata(pages, metadata):
    from itertools import product
    for page, meta in product(pages, metadata):
        if page.id == meta.id:
            page.title = meta.title
            page.description = meta.description
    return pages

    # Access the database abd insert the stuff


