from api.helpers.verification import get_current_active_user
from uuid import uuid4
from fastapi import APIRouter
from fastapi.param_functions import Depends
from api.db.connection import database, pages
from api.helpers.models import Page, User, UserOut
from api.helpers.pages import unwrap_submitted_page

router = APIRouter(
    prefix="/pages",
    tags=["Pages"]
)

@router.post('/add_page', response_model=Page)
async def save_new_page(new_page: Page):
    new_page.id = uuid4()
    query = pages.insert().values(**new_page.dict())
    await database.execute(query=query)
    return new_page

@router.get('/protected')
async def protected_page(new_page: Page = Depends(unwrap_submitted_page), current_user: User = Depends(get_current_active_user)):
    current_user = UserOut(**current_user.dict())
    return {"message": "Protected URL accessed by user.", "url": new_page, **current_user.dict()}

