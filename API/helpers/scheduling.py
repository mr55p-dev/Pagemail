import async_scheduler as asch
from apscheduler.jobstores.memory import MemoryJobStore
from apscheduler.schedulers.asyncio import AsyncIOScheduler

# Create a scheduler for the metadata update job.
memory_job = MemoryJobStore()
scheduler = AsyncIOScheduler()
scheduler.add_jobstore(memory_job, alias="local")

my_scheduler = asch.Scheduler(None, target_func=lambda x: x)