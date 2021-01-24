from uuid import uuid4
from fastapi import APIRouter
from api.db.connection import database, pages
from api.db.models import SavePage

router = APIRouter(
    prefix="/v1/pages",
    tags=["Pages"]
)

@router.post('/add_page', response_model=SavePage)
async def save_new_page(new_page: SavePage):
    new_page.id = uuid4()
    query = pages.insert().values(**new_page.dict())
    await database.execute(query=query)
    return new_page

