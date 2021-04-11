import logging

from datetime import datetime, timedelta
from uuid import uuid4

from async_scheduler.job import UserJob
from async_scheduler.exceptions import DuplicateJobError

from fastapi import APIRouter, Depends, Form, HTTPException
from fastapi.security.oauth2 import OAuth2PasswordRequestForm

from API.db.operations import user_create, user_delete, user_read
from API.helpers.email_tools import send_email
from API.helpers.models import BaseEmail, UserIn, UserOut
from API.helpers.scheduling import sch_news
from API.helpers.verification import (create_new_token, get_validated_user,
                                      hash_password, validate_user)


router = APIRouter(
    prefix="/user",
    tags=["Users"]
)

log = logging.getLogger(__name__)
log.setLevel(logging.DEBUG)

def decode_new_user_form(email: str = Form(...), name: str = Form(...), password: str = Form(...)):
    """Create a new user from a form"""
    # Catch validation error of the email/password here.
    return UserIn(name=name, email=email, password=password)

def decode_user_form(email: str = Form(...), password: str = Form(...)):
    """Create a user object from form submission"""
    return UserIn(email=email, password=password)

# , response_model=UserOut <= Change it to a dictionary of an api key response and user.
@router.post('/register')
async def add_user(new_user: UserIn = Depends(decode_new_user_form)):
    """Add a new user to the table, send the welcome email and generate a token"""
    new_user.id = uuid4()
    new_user.password = hash_password(new_user.password)
    new_user.date_added = datetime.utcnow()
    new_user.is_active = True

    await user_create(new_user)

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

@router.delete('/self', response_model=UserOut)
async def delete_user(to_delete: UserIn = Depends(decode_user_form)):
    """Delete a user from the table (requires username/password auth)"""
    user: UserIn = await user_read(to_delete.email)
    validate_user(to_delete.password, user.password)
    await user_delete(user)
    return UserOut(**user.dict())

@router.get('/self', response_model=UserOut)
async def read_self(current_user: UserIn = Depends(get_validated_user)):
    """Get a users information from table"""
    api_token = create_new_token(
        {"sub": current_user.email},
        page_only=False)
    return UserOut(**current_user.dict(), token=api_token)

# @router.put('/self')
# async def modify_user(
#     current_user: UserIn = Depends(decode_user_form),
#     new_fields = Depends(...)):
#     """Update a users information in the table"""
#     user = await user_read(current_user.email)
#     log.debug(new_fields)
#     validate_user(current_user.password, user.password)
    # Make sure that the ID cannot be updated.
    # Update the user
    # await user_update()
    # return 200

@router.post('/token')
async def login(form_data: OAuth2PasswordRequestForm = Depends(), page_only: int = Form(0)):
    """Fetch a token based on login details"""
    user = await user_read(form_data.username)
    validate_user(form_data.password, user.password)
    token = create_new_token(
        {"sub": user.email},
        page_only=page_only
    )
    return {"access_token": token, "token_type": "bearer", "user": UserOut(**user.dict())}

@router.post('/newsletter')
async def test_send_email(current_user: UserIn = Depends(get_validated_user)):
    """Subscribe a user to the newsletter"""
    # message = BaseEmail(
    #     subject="Test mail!",
    #     content="This is a test email. Please disregard it.",
    #     recipients=current_user.email
    # )
    # await send_email(message)
    # message.content = "This is a test email sent with the scheduler. Please disregard it."
    # # scheduler.add_job(send_email, args=[message])
    # return 200
    job = UserJob(timedelta(seconds=15), current_user.id)
    try:
        sch_news.add(job)
    except DuplicateJobError as exception:
        raise HTTPException(
            status_code=400,
            detail="the user is already subscribed to emails."
        ) from exception
    return 200

@router.delete('/newsletter')
async def dont_email_me(current_user: UserIn = Depends(get_validated_user)):
    """Unsubscribe a user from the newsletter"""
    try:
        job = sch_news.pop_user(current_user.id)
    except ValueError as exception:
        raise HTTPException(
            status_code=400,
            detail="the user is not subscribed to emails."
        ) from exception
    return job
