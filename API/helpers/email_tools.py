import logging
import os
from email.message import EmailMessage
from uuid import UUID

import aiosmtplib as sm
from API.helpers.models import BaseEmail

# Create an application level ssl context
config = {
    "hostname": os.getenv('MAIL_HOSTNAME'),
    "username": os.getenv('MAIL_USERNAME'),
    "password": os.getenv('MAIL_PASSWORD'),
    "timeout" : int(os.getenv('MAIL_TIMEOUT')),
    "use_tls": True
}

async def send_email(mail: BaseEmail = None):
    """Configure and send an email over SSL
    Args:
        mail (BaseEmail): The recipient(s), subject and contents of the message."""
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

async def newsletter(user_id: UUID = None, job_type: str = None, **kwargs):
    """Newsletter sending function (to be used with a scheduler"""
    del kwargs
    log = logging.getLogger(__name__)
    log.setLevel(logging.DEBUG)

    log.debug("Sending newsletter for user: %s with job_type: %s", {user_id}, {job_type})
    # Fetch the user
    # Compose the email:
    # - Refactor the "current user" and "pages" stuff into db/requests.py
    # send the email!
    return True
