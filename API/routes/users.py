from datetime import datetime
from uuid import uuid4

from API.db.connection import get_db, page_metadata, pages, users
from API.helpers.email_tools import send_email
from API.helpers.scheduling import my_scheduler
from API.helpers.models import BaseEmail, UserIn, UserOut
from API.helpers.verification import (create_new_token, fetch_user,
                                      get_validated_user, hash_password,
                                      validate_user)
from async_scheduler.job import Job
from asyncpg.exceptions import UniqueViolationError
from databases.core import Database
from fastapi import APIRouter, Depends, Form, HTTPException
from fastapi.security.oauth2 import OAuth2PasswordRequestForm
from sqlalchemy import select

router = APIRouter(
    prefix="/user",
    tags=["Users"]
)

def decode_new_user_form(email: str = Form(...), name: str = Form(...), password: str = Form(...)):
    return UserIn(name=name, email=email, password=password)

def decode_user_form(email: str = Form(...), password: str = Form(...)):
    return UserIn(email=email, password=password)

# , response_model=UserOut <= Change it to a dictionary of an api key response and user.
@router.post('/register')
async def add_user(new_user: UserIn = Depends(decode_new_user_form),
        database: Database = Depends(get_db)):

    new_user.id = uuid4()
    new_user.password = hash_password(new_user.password)
    new_user.date_added = datetime.utcnow()
    new_user.is_active = True

    query = users.insert().values(**new_user.dict())
    try:
        await database.execute(query=query)
    except UniqueViolationError:
        raise HTTPException(
            status_code=400,
            detail="Username already exists."
        )
    # ^Factorise this into a helper function
    # Queue the email
    onboarding_email = BaseEmail(
        recipients=new_user.email,
        subject="Welcome to Pagemail!",
        content=f"Hello {new_user.name}, and welcome to Pagemail!"
    )
    send_email(**{"mail": onboarding_email})
    # Collect a token
    token = create_new_token({"sub": new_user.email})
    return {
        "user": UserOut(**new_user.dict()),
        "token": {
            "access_token": token,
            "token_type": "bearer"
        }
    }

@router.delete('/remove', response_model=UserOut)
async def delete_user(to_delete: UserIn = Depends(decode_user_form),
        database: Database = Depends(get_db)):

    user = await fetch_user(to_delete.email)
    validate_user(to_delete.password, user.password)
    query = select([pages.c.id]).select_from(users.join(pages)).where(users.c.id == user.id)
    pages_owned = await database.fetch_all(query)
    page_ids = [i["id"] for i in pages_owned]

    query1 = page_metadata.delete().where(page_metadata.c.id.in_(page_ids))
    query2 = pages.delete().where(pages.c.id.in_(page_ids))
    query3 = users.delete().where(users.c.id == user.id)

    await database.execute(query1)
    await database.execute(query2)
    await database.execute(query3)

    return UserOut(**user.dict())

@router.get('/self', response_model=UserOut)
async def read_self(current_user: UserIn = Depends(get_validated_user)):
    api_token = create_new_token(
        {
            "sub": current_user.email
        },
        page_only=True)
    current_user = UserOut(**current_user.dict(), token=api_token)
    return current_user

@router.post('/token')
async def login(form_data: OAuth2PasswordRequestForm = Depends(), page_only: int = Form(0)):
    user = await fetch_user(form_data.username)
    validate_user(form_data.password, user.password)
    token = create_new_token(
        {"sub": user.email},
        page_only=page_only
    )
    return {"access_token": token, "token_type": "bearer", "user": UserOut(**user.dict())}

@router.get('/subscribe')
async def test_send_email(current_user: UserIn = Depends(get_validated_user), db = Depends(get_db)):
    from datetime import timedelta
    message = BaseEmail(
        subject="Test mail!",
        content="This is a test email. Please disregard it.",
        recipients=current_user.email
    )
    # await send_email(message)
    # message.content = "This is a test email sent with the scheduler. Please disregard it."
    # # scheduler.add_job(send_email, args=[message])
    # return 200
    job = Job(current_user.id, timedelta(seconds=15))
    try:
        my_scheduler.add(job)
    except:
        raise HTTPException(
            status_code=400,
            detail="the user is already subscribed to emails."
        )
    return 200

@router.get('/unsubscribe')
async def dont_email_me(current_user: UserIn = Depends(get_validated_user)):
    try:
        job = my_scheduler.pop_user(current_user.id)
    except ValueError:
        raise HTTPException(
            status=400,
            detail="the user is not subscribed to emails."
        )
    return job
