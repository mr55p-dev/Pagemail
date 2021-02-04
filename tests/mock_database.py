import asyncio
import os

import databases
from API.db.connection import metadata, pages, users
from sqlalchemy import create_engine

DATABASE_URL = os.getenv("TESTING_DATABASE_URL")
print(DATABASE_URL)

database = databases.Database(DATABASE_URL)
engine = create_engine(DATABASE_URL)
metadata.bind = engine
metadata.create_all(engine)

async def reset_tables():
    from tests.data_gen import get_mock_data

    print("Dropping all tables...")
    metadata.drop_all(engine)
    print("Done")
    print("Recreating tables...")
    metadata.create_all(engine)
    print("Done")

    print("Populating user table...")
    (mock_users, mock_pages) = get_mock_data()
    query = users.insert()
    await database.execute_many(query=query, values=mock_users)
    print("Done")

    print("Populating pages table...")
    query = pages.insert()
    await database.execute_many(query=query, values=mock_pages)
    print("Done")

loop = asyncio.get_event_loop()
loop.run_until_complete(database.connect())
loop.run_until_complete(reset_tables())

def override_get_db():
    yield database
