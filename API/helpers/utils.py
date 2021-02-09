from uuid import UUID
from API.helpers.models import Page
from fastapi import Form, HTTPException
from pydantic.error_wrappers import ValidationError
import requests
from bs4 import BeautifulSoup as BS


async def unwrap_submitted_page(url: str = Form("")):
    try:
        page = Page(page_url=url)
    except ValidationError:
        raise HTTPException(
        status_code=400,
        detail="Malformed request.",
    )
    return page

async def fetch_metadata(id: UUID, url: str):
    resp = requests.get(url)
    # Validate response

    # Get page title
    soup = BS(resp.text)
    title = soup.title.text

    # 


