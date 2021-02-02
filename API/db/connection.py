from sqlalchemy import *
from sqlalchemy.sql.functions import now
from sqlalchemy.dialects.postgresql import UUID
from uuid import uuid4
from dotenv import load_dotenv
import databases
import os

# Get environment variables
DATABASE_URL = os.getenv("DATABASE_URL")

metadata = MetaData()

# Define tables
users = Table(
    "users",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, default=uuid4, unique=True),
    Column("name", String, nullable=False, unique=True),
    Column("email", String, nullable=False),
    Column("password", String, nullable=False),
    Column("date_added", DateTime, server_default=now(), nullable=False),
    Column("is_active", Boolean, nullable=False),
    )

pages = Table(
    "pages",
    metadata,
    Column("id", UUID(as_uuid=False), primary_key=True, unique=True, default=uuid4),
    Column("page_url", String, nullable=False),
    Column("date_added", DateTime, server_default=now(), nullable=False),
    Column("user_id", UUID(as_uuid=False), ForeignKey("users.id"))
)

jobs = Table(
    "jobs", metadata,
    Column('id', Unicode(191, _warn_on_bytestring=False), primary_key=True),
    Column('next_run_time', Float(25), index=True),
    Column('job_state', LargeBinary, nullable=False),
    Column('user_id', UUID(as_uuid=False), ForeignKey('users.id'))
)

# Create database connection and tables.
database = databases.Database(DATABASE_URL)
engine = create_engine(DATABASE_URL)
metadata.create_all(engine)