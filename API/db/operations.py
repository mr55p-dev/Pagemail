from typing import List, Optional, Union
from uuid import UUID

from itertools import product
from asyncpg.exceptions import UniqueViolationError
from databases.core import Database
from fastapi import HTTPException, status
from pydantic.error_wrappers import ValidationError
from sqlalchemy import select
from starlette.status import HTTP_400_BAD_REQUEST
from API.db.connection import get_db, page_metadata, pages, users
from API.helpers.models import (Page, PageFilled, PageMetadata, PageOut,
                                UserIn, UserOut)

def db_connection() -> Database:
    """Returns an instance of the database connection"""
    return next(get_db())

def set_page_metadata(pages_list, metadata):
    """Sets the metadata of a given list of pages"""
    for page, meta in product(pages_list, metadata):
        if page.id == meta.id:
            page.title = meta.title
            page.description = meta.description
    return pages_list


async def meta_create(meta: PageMetadata):
    """Create a metadata entry
    Args:
        meta (PageMetadata): The metadata object to insert
    Returns:
        DatabaseRow: Some kind of generator object."""
    database = db_connection()
    query = page_metadata.insert(None).values(**meta.dict())
    response = await database.execute(query)
    return response

async def meta_read(metadata_id: UUID) -> PageMetadata:
    """Reads a row from the metadata table
    Args:
        metadata_id (UUID): The id of the metadata to fetch
    Returns:
        PageMetadata"""
    database = db_connection()
    query = page_metadata.select().where(page_metadata.c.id == metadata_id)
    response = database.fetch_one(query)
    return PageMetadata(**response)

async def metas_read(page_ids: List[UUID]) -> List[PageMetadata]:
    """Reads multiple rows from the metadata table
    Args:
        page_ids (List[UUID]): The referenced page ids
    Returns:
        List[MetaData]"""
    database = db_connection()
    query = page_metadata.select().where(page_metadata.c.id.in_(page_ids))
    metadata = await database.fetch_all(query)
    return [PageMetadata(**i) for i in metadata]


async def user_create(new_user: UserIn) -> bool:
    """Create a new user and insert into the database
    Args:
        new_user (UserIn): The user to insert
    Returns:
        : True if sucessful
    """
    database = db_connection()
    query = users.insert(None).values(**new_user.dict())
    try:
        response = await database.execute(query)
    except UniqueViolationError as exception:
        raise  HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="This user already exists."
        ) from exception
    except ValidationError as exception:
        raise HTTPException(
            status_code=HTTP_400_BAD_REQUEST,
            detail="This is an invalid email address or password"
        ) from exception
    return UserOut(**new_user.dict()) if response else None

async def user_read(user_email: str = None, user_id: UUID = None) -> UserIn:
    """Read a user from the table
    Args:
        user_email (str): The users email
        user_id (UUID): The users id
    Returns:
        user (UserOut): UserOut object with the users info
    """
    database = db_connection()
    if not (user_email or user_id):
        return None
    elif not user_id:
        query = users.select().where(users.c.email == user_email)
    else:
        query = users.select().where(users.c.id == user_id)

    response = await database.fetch_one(query)

    if not response:
        raise HTTPException(
            status_code=400,
            detail="Incorrect user information."
            )
    return UserIn(**response)

async def user_update(updated_user: UserIn) -> UserOut:
    """Update a user in the table
    Args:
        updated_user (UserIn): User to update (must contain id field)
    Returns:
        user_out (UserIn): The updated user
    """
    database = db_connection()
    query = users.update(None).where(users.c.id == updated_user.id).values(**updated_user.dict())
    response = await database.execute(query)
    return UserOut(**response)

async def user_delete(user: UserIn) -> UserOut:
    """Delete a user from the table
    Args:
        user (UserIn): User to delete
    Returns:
        user (UserOut): Deleted user
    """
    database = db_connection()

    # Fetch the pages owned by the user
    query = (
                select([pages.c.id])
                .select_from(users.join(pages))
                .where(users.c.id == user.id)
            )
    pages_owned = await database.fetch_all(query)
    page_ids = [i["id"] for i in pages_owned]

    # Delete the metadata
    query = page_metadata.delete(None).where(page_metadata.c.id.in_(page_ids))
    await database.execute(query)

    # Delete the pages
    query = pages.delete(None).where(pages.c.id.in_(page_ids))
    await database.execute(query)

    # Delete the user
    query = users.delete(None).where(users.c.id == user.id)
    await database.execute(query)

    return UserOut(**user.dict())


async def page_create(new_page: Page) -> PageOut:
    """Create a new page in the table
    Args:
        new_page (Page): The page to insert
    Returns:
        page (Page): The inserted page
    """
    database = db_connection()
    query = pages.insert(None).values(**new_page.dict())
    response = await database.execute(query)
    return PageOut(**response)

async def page_read(page_id: UUID = None) -> Optional[Union[PageFilled, List[PageOut]]]:
    """Read a page or pages from the table using either the user or page id
    Args:
        user_id (UUID): User id to fetch all pages for
        page_id (UUID): An individual page to fetch
    Returns:
        None: If neither argument is specified
        Page (PageOut): A single page if an id is specified
        Pages (List[PageOut]): A list of pages owned by a user
    Raises:
        HTTPException (404): The page requested does not exist
    """
    database = db_connection()
    query = pages.select().where(pages.c.id == page_id)
    response = await database.fetch_one(query)
    if not response:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="The requested page does not exist."
        )
    meta = await meta_read((page := PageFilled(**response)).id)
    page = set_page_metadata(page, meta) if meta else page
    return page


async def pages_read(user_id: UUID) -> List[PageFilled]:
    """Reads pages from the table
    Args:
        user_id (UUID): User id for which to select the pages.
    Returns:
        List[PageFilled]: Page objects populated with metadata"""
    database = db_connection()
    query = pages.select().where(pages.c.user_id == user_id)
    user_pages = await database.fetch_all(query)
    user_pages = [PageFilled(**i) for i in user_pages]

    page_ids = [i.id for i in user_pages]
    page_meta = await metas_read(page_ids)

    return set_page_metadata(user_pages, page_meta)


async def page_update(updated_page: Page) -> PageOut:
    """Update a page in the table
    Args:
        updated_page (Page): The page with updated values
    Returns:
        PageOut: The updated page.
    """
    database = db_connection()
    query = (
                pages
                .update(None)
                .where(pages.c.id == updated_page.id)
                .values(**updated_page.dict())
            )
    await database.execute(query)
    return PageOut(**updated_page.dict())

async def page_delete(deleted_page: Page) -> PageOut:
    """Deletes a page from the table
    Args:
        deleted_page (Page): The page to delete
    Returns:
        PageOut: The deleted page
    """
    database = db_connection()

    # Delete the metadata of the page
    # Change to delete_meta in the future
    query = page_metadata.delete(None).where(page_metadata.c.id == deleted_page.id)
    await database.execute(query)

    # Delete the page
    query = pages.delete(None).where(pages.c.id == deleted_page.id)
    await database.execute(query)
    return PageOut(**deleted_page.dict())


async def page_verify(user: UserIn, page_id: UUID) -> Optional[Page]:
    """Verify that the given user owns the page in question
    Args:
        user (UserIn): The user making the request
        page_id (UUID): The page in question
    Returns:
        False: The page is not owned by the user
        Page: The page is owned by the user
    """
    page = await page_read(page_id=page_id)
    if page.user_id == user.id:
        return Page(**page)
    else:
        return False

async def pages_no_meta() -> Optional[List[Page]]:
    """Fetch all the pages for which there are no metadata
    Returns:
        List[Page]: A list of pages
        None: If there are no pages
    """
    database = db_connection()
    query = pages.join(page_metadata, full=True)
    query = select([pages]).select_from(query).where(page_metadata.c.id is None)
    response = await database.fetch_all(query)
    return response
