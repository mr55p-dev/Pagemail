from email.message import EmailMessage
import aiosmtplib as sm
import os
import ssl
from API.helpers.models import BaseEmail

# Create an application level ssl context
ssl_context = ssl.create_default_context()
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
