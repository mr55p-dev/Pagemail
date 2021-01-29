from api.helpers.verification import get_current_active_user
from uuid import uuid4
from fastapi import APIRouter
from fastapi.param_functions import Depends
from api.db.connection import database, pages
from api.helpers.models import Page, User, UserOut
from api.helpers.pages import unwrap_submitted_page
from datetime import datetime

router = APIRouter(
    prefix="/page",
    tags=["Pages"]
)

@router.post('/add_page', response_model=Page)
async def save_new_page(
    new_page: Page = Depends(unwrap_submitted_page),
    current_user: User = Depends(get_current_active_user)
    ):
    new_page.id = uuid4()
    new_page.user_id = current_user.id
    new_page.date_added = datetime.now()

    query = pages.insert().values(**new_page.dict())
    await database.execute(query=query)
    return new_page

