from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.jobstores.sqlalchemy import SQLAlchemyJobStore

from API.db.connection import engine, jobs

# Use a custom table design which can store an associated user_id along with the
# normal table id.
class SQLAlchemyJob(SQLAlchemyJobStore):
    def override_table(self, table):
        self.jobs_t = table


job_store = SQLAlchemyJob(engine=engine, tablename="dummy_jobs")
job_store.override_table(jobs)
# Recreate the table however with ID as a foreign key to the users table.

scheduler = AsyncIOScheduler()
scheduler.add_jobstore(job_store)