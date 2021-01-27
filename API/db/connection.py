from sqlalchemy import *
from sqlalchemy.sql.functions import now
from sqlalchemy.dialects.postgresql import UUID
from uuid import uuid4
from dotenv import load_dotenv
import databases
import os

# Get environment variables
DATABASE_URL = os.getenv("DATABASE_URI")


metadata = MetaData()

# Define tables
users = Table(
    "users",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, default=uuid4, unique=True),
    Column("name", String),
    Column("email", String),
    Column("password", String),
    Column("date_added", DateTime, server_default=now()),
    Column("is_active", Boolean, default=True)
    )

pages = Table(
    "pages",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, unique=True),
    Column("page_url", String),
    Column("date_added", DateTime, server_default=now())
)

# Create database connection and tables.
database = databases.Database(DATABASE_URL)
engine = create_engine(DATABASE_URL)
metadata.create_all(engine)