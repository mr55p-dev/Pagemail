from sqlalchemy import *
from sqlalchemy.sql.functions import now
from sqlalchemy.dialects.postgresql import UUID
from uuid import uuid4
from dotenv import load_dotenv
import databases
import os

# Get environment variables
load_dotenv(verbose=True)
DATABASE_URL=os.getenv("DATABASE_URI")
SECRET = os.getenv("SECRET_KEY")


metadata = MetaData()

# Define tables
users = Table(
    "users",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, default=uuid4, unique=True),
    Column("name", String),
    Column("email", String),
    Column("password_hash", String),
    Column("date_added", DateTime, server_default=now()),
    )

pages = Table(
    "pages",
    metadata,
    Column("id", UUID(as_uuid=True), primary_key=True, unique=True),
    Column("page_url", String),
    Column("date_added", DateTime, server_default=now())
)

# Create database connection and tables.
database = databases.Database(DATABASE_URL)
engine = create_engine(DATABASE_URL)
metadata.create_all(engine)