import async_scheduler as asch
from apscheduler.jobstores.memory import MemoryJobStore
from apscheduler.schedulers.asyncio import AsyncIOScheduler

# Create a scheduler for the metadata update job.
memory_job = MemoryJobStore()
scheduler = AsyncIOScheduler()
scheduler.add_jobstore(memory_job, alias="local")

sch_news = asch.Scheduler()
sch_meta = asch.Scheduler()