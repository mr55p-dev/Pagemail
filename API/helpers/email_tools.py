from email.message import EmailMessage
import aiosmtplib as sm
import os
import ssl
from API.helpers.models import NewsletterEmail

# Create an application level ssl context
ssl_context = ssl.create_default_context()
config = {
    "hostname": os.getenv('MAIL_HOSTNAME'),
    "username": os.getenv('MAIL_USERNAME'),
    "password": os.getenv('MAIL_PASSWORD'),
    "timeout" : os.getenv('MAIL_TIMEOUT')
}

async def send_email(mail: NewsletterEmail):
    # Set up an EmailMessage object
    message = EmailMessage()
    message["From"] = mail.sender
    message["To"] = mail.recipients
    message["Subject"] = mail.subject
    message.set_content = mail.content

    # Send the email with configuration settings.
    await sm.send(
        message=message,
        **config
    )

# async def send_test_email():
#     # Setup the message
#     message = EmailMessage()
#     message["From"] = "elunnon.pagemail@gmail.com"
#     message["To"] = "ellislunnon@gmail.com"
#     message["Subject"] = "Hello World"
#     message.set_content = "Hello, this is an automated email."
#     # Look into creating a connection to SMTP on startup, or batches of connections to reduce overhead.
#     await sm.send(
#         message=message,
#         hostname="smtp.gmail.com",
#         use_tls=True,
#         tls_context=ssl_context,
#         username=os.getenv('MAIL_USERNAME'),
#         password=os.getenv('MAIL_PASSWORD'),
#         timeout=30
#         )
