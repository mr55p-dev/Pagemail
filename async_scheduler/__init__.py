"""Implement:
[ ] Proper naming of the jobs
[x] Self updating NEXT_CALL_DELAY
[ ] Change from having one target function to storing a dict of target functions for different job types/roles
[ ] Save the jobs
[ ] Database stuff (should it be part of a separate module?
    Could be enabled/disabled/easier to modify)
"""

from async_scheduler.job import Job
from async_scheduler.scheduler import Scheduler

__version__ = "0.0.1"

