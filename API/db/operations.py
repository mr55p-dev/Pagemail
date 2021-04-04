from logging import raiseExceptions
from typing import List, Optional, Union
from uuid import UUID

from sqlalchemy.sql.expression import delete
from API.db.connection import get_db, page_metadata, pages, users
from API.helpers.models import Page, PageOut, UserIn, UserOut
from asyncpg.exceptions import UniqueViolationError
from databases.core import Database
from fastapi import Depends, HTTPException, status
from fastapi.exceptions import HTTPException
from starlette.status import HTTP_400_BAD_REQUEST
from sqlalchemy import select

# async def execute(query, error_callback: callable = base_error_callback):
#     database = next(get_db())
#     try:
#         response = database.execute(query)
#     except Exception as e:
#         error_callback(e)
#     return response or None


async def user_create(new_user: UserIn,) -> bool:
    """Create a new user and insert into the database
    Args:
        new_user (UserIn): The user to insert
    Returns:
        return_code (bool): True if sucessful
    """
    database = next(get_db())
    query = users.insert().values(**new_user.dict())
    try:
        response = await database.execute(query)
    except UniqueViolationError:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="This user already exists."
        )
    return UserOut(**response)


async def user_read(user_email: str = None, user_id: UUID = None) -> UserIn:
    """Read a user from the table
    Args:
        user_email (str): The users email
        user_id (UUID): The users id
    Returns:
        user (UserOut): UserOut object with the users info
    """
    database = next(get_db())
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
    return UserOut(**response)

async def user_update(updated_user: UserIn) -> UserOut:
    """Update a user in the table
    Args:
        updated_user (UserIn): User to update (must contain id field)
    Returns:
        user_out (UserIn): The updated user
    """
    database = next(get_db())
    query = users.update().where(users.c.id == updated_user.id).values(**updated_user.dict())
    response = await database.execute(query)
    return UserOut(**response)

async def user_delete(user: UserIn) -> UserOut:
    """Delete a user from the table
    Args:
        user (UserIn): User to delete
    Returns:
        user (UserOut): Deleted user
    """
    database = next(get_db())

    # Fetch the pages owned by the user
    query = (
                select([pages.c.id])
                .select_from(users.join(pages))
                .where(users.c.id == user.id)
            )
    pages_owned = await database.fetch_all(query)
    page_ids = [i["id"] for i in pages_owned]

    # Delete the metadata
    query = page_metadata.delete().where(page_metadata.c.id.in_(page_ids))
    await database.execute(query)

    # Delete the pages
    query = pages.delete().where(pages.c.id.in_(page_ids))
    await database.execute(query)

    # Delete the user
    query = users.delete().where(users.c.id == user.id)
    await database.execute(query)

    return UserOut(**user.dict())


async def page_create(new_page: Page) -> PageOut:
    """Create a new page in the table
    Args:
        new_page (Page): The page to insert
    Returns:
        page (Page): The inserted page
    """
    database = next(get_db())
    query = pages.insert().values(**new_page.dict())
    response = await database.execute(query)
    return PageOut(**response)

async def read_page(user_id: UUID = None, page_id: UUID = None) -> Optional[Union[PageOut, List[PageOut]]]:
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
    database = next(get_db())
    if not (user_id or page_id):
        return None
    elif user_id:
        query = pages.select().where(pages.c.user_id == user_id)
        response = await database.fetch_all(query)
        return [PageOut(**i) for i in response]
    else:
        query = pages.select().where(pages.c.id == page_id)
        response = await database.fetch_one(query)
        if not response:
            raise HTTPException(
                status=status.HTTP_404_NOT_FOUND,
                detail="The requested page does not exist."
            )
        return PageOut(**response)

async def update_page(updated_page: Page) -> PageOut:
    """Update a page in the table
    Args:
        updated_page (Page): The page with updated values
    Returns:
        PageOut: The updated page.
    """
    database = next(get_db())
    query = (
                pages
                .update()
                .where(pages.c.id == updated_page.id)
                .values(**updated_page.dict())
            )
    response = await database.execute(query)
    return PageOut(**updated_page.dict())

async def delete_page(deleted_page: Page) -> PageOut:
    """Deletes a page from the table
    Args:
        deleted_page (Page): The page to delete
    Returns:
        PageOut: The deleted page
    """
    database = next(get_db())

    # Delete the metadata of the page
    query = page_metadata.delete().where(page_metadata.c.id == deleted_page.id)
    await database.execute(query)

    # Delete the page
    query = pages.delete().where(pages.c.id == deleted_page.id)
    await database.execute(query)
    return PageOut(**deleted_page.dict())

async def verify_page(user: UserIn, page_id: UUID) -> Optional[Page]:
    """Verify that the given user owns the page in question
    Args:
        user (UserIn): The user making the request
        page_id (UUID): The page in question
    Returns:
        False: The page is not owned by the user
        Page: The page is owned by the user
    """
    page = await read_page(page_id=page_id)
    if page.user_id == user.id:
        return Page(**page)
    else:
        return False

async def without_metadata_pages() -> Optional[List[Page]]:
    """Fetch all the pages for which there are no metadata
    Returns:
        List[Page]: A list of pages
        None: If there are no pages
    """
    database = next(get_db())
    query = pages.join(page_metadata, full=True)
    query = select([pages]).select_from(query).where(page_metadata.c.id == None)
    response = await database.fetch_all(query)
    return response





