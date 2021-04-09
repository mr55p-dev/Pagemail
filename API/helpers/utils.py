import logging
from uuid import UUID

import requests
from bs4 import BeautifulSoup as BS
from requests.exceptions import RequestException
from API.db.operations import meta_create, page_read, pages_no_meta, user_read
from API.helpers.models import Page, PageMetadata, UserIn

log = logging.getLogger(__name__)
log.setLevel(logging.DEBUG)

async def update_metadata():
    """Scheduled action to update the metadata in the index"""
    unannotated_pages = await pages_no_meta()
    for page in unannotated_pages:
        await fetch_metadata(page_id=page.id, url=page.url)
    log.debug("%d page(s) have had their metadata updated.", len(unannotated_pages))
    return unannotated_pages

async def verify_ownership(user_id, page_id) -> bool:
    """Verify that the given user_id owns the page_id
    Args:
        user_id (UUID)
        page_id (UUID)
    Returns:
        bool"""
    user: UserIn = await user_read(user_id=user_id)
    page: Page = await page_read(page_id=page_id)
    return page.user_id == user.id

def fetch_page_information(url: str):
    """Fetches a page and extracts the title and description.
    Args:
        url (str): The url of the site to scrape
    Returns:
        title, description (tuple): The fetched page information"""
    with requests.Session() as session:
        session.headers = {
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0",
            "Accept": "*/*",
            "Accept-Encoding": "gzip, deflate, br",
            "Accept-Language": "en-GB,en;q=0.5",
            "Connection": "close",
            "DNT": 1
        },

    try:
        resp = session.get(url)
    except RequestException:
        title = ""
        desc = ""
    else:
        soup = BS(resp.text, features="html.parser")
        title = soup.title.text
        desc = soup.get_text()[len(title):500]
    return (title, desc)


async def fetch_metadata(page_id: UUID, url: str) -> None:
    """Worker task to fetch the metadata of a page
    to keep as a description of it"""
    title, description = fetch_page_information(url)
    meta = PageMetadata(
        id=page_id,
        title=title,
        description=description
    )
    await meta_create(meta)
