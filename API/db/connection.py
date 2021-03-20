from sqlalchemy import *
from sqlalchemy.sql.functions import now
from sqlalchemy.dialects.postgresql import UUID
from uuid import uuid4
from dotenv import load_dotenv
import databases
import os
import ssl

# Set SSL options to enable postgres connection
context = ssl.create_default_context()
context.check_hostname = False
context.verify_mode = ssl.CERT_NONE


# Get environment variables
DATABASE_URL = os.getenv("DATABASE_URL")
DATABASE_USE_SSL = False if os.getenv("DATABASE_USE_SSL") == "0" else context
DATABASE_MIN_CONNECTIONS = int(os.getenv("DATABASE_MIN_CONNECTIONS"))
DATABASE_MAX_CONNECTIONS = int(os.getenv("DATABASE_MAX_CONNECTIONS"))


metadata = MetaData()

# Define tables
users = Table(
    "users",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, default=uuid4, unique=True),
    Column("name", String, nullable=False, unique=False),
    Column("email", String, nullable=False, unique=True),
    Column("password", String, nullable=False),
    Column("date_added", DateTime, server_default=now(), nullable=False),
    Column("is_active", Boolean, nullable=False),
    )

pages = Table(
    "pages",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, unique=True, default=uuid4),
    Column("url", String, nullable=False),
    Column("date_added", DateTime, server_default=now(), nullable=False),
    Column("user_id", UUID(as_uuid=False), ForeignKey("users.id"))
)

page_metadata = Table(
    'page_metadata',
    metadata,
    Column('id', UUID(as_uuid=False), ForeignKey("pages.id"), primary_key=True, unique=True),
    Column('title', String),
    Column('description', String)
)

jobs = Table(
    "jobs", metadata,
    Column('id', Unicode(191, _warn_on_bytestring=False), primary_key=True),
    Column('next_run_time', Float(25), index=True),
    Column('job_state', LargeBinary, nullable=False),
    Column('user_id', UUID(as_uuid=False), ForeignKey('users.id'))
)

# Create database connection and tables.
database = databases.Database(DATABASE_URL, ssl=DATABASE_USE_SSL, min_size=DATABASE_MIN_CONNECTIONS, max_size=DATABASE_MAX_CONNECTIONS)
engine = create_engine(DATABASE_URL)
metadata.create_all(engine)

def get_db():
    yield database