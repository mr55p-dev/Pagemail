from email.message import EmailMessage
import aiosmtplib as sm
import os
from API.helpers.models import BaseEmail
from API.db.connection import context as ssl_context

# Create an application level ssl context
config = {
    "hostname": os.getenv('MAIL_HOSTNAME'),
    "username": os.getenv('MAIL_USERNAME'),
    "password": os.getenv('MAIL_PASSWORD'),
    "timeout" : int(os.getenv('MAIL_TIMEOUT')),
    "use_tls": True
}

async def send_email(mail: BaseEmail):
    # Set up an EmailMessage object
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
