import logging
import os
from email.message import EmailMessage
from uuid import UUID

import aiosmtplib as sm
from API.db.connection import context as ssl_context
from API.helpers.models import BaseEmail

# Create an application level ssl context
config = {
    "hostname": os.getenv('MAIL_HOSTNAME'),
    "username": os.getenv('MAIL_USERNAME'),
    "password": os.getenv('MAIL_PASSWORD'),
    "timeout" : int(os.getenv('MAIL_TIMEOUT')),
    "use_tls": True
}

async def send_email(mail: BaseEmail = None, *args, **kwargs):
    # Set up an EmailMessage object
    if not mail:
        return None
    message = EmailMessage()
    message["From"] = mail.sender
    message["To"] = mail.recipients
    message["Subject"] = mail.subject
    message.set_content(mail.content)

    # Send the email with configuration settings.
    await sm.send(
        message=message,
        **config
    )

async def newsletter(*args, user_id: UUID = None, job_type: str = None):
    log = logging.getLogger(__name__)
    log.setLevel(logging.DEBUG)

    log.debug(f"Sending newsletter for user: {user_id} with job_type: {job_type}")
    # Fetch the user
    # Compose the email:
    # - Refactor the "current user" and "pages" stuff into db/requests.py
    # send the email!
    return True