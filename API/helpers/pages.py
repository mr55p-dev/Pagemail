from api.helpers.models import Page
from fastapi import Form, HTTPException
from pydantic.error_wrappers import ValidationError


async def unwrap_submitted_page(url: str = Form("")):
    try:
        page = Page(page_url=url)
    except ValidationError:
        raise HTTPException(
        status_code=400,
        detail="Malformed request.",
    )
    return page