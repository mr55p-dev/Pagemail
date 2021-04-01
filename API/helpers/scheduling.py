from asyncio.futures import Future
import functools
from types import coroutine
from typing import Any, Awaitable, Iterable, List
from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.jobstores.memory import MemoryJobStore

# Create a scheduler for the metadata update job.
memory_job = MemoryJobStore()
scheduler = AsyncIOScheduler()
scheduler.add_jobstore(memory_job, alias="local")
# %%
import logging
import asyncio
from functools import partial
from abc import abstractmethod
from datetime import datetime as dt
from uuid import UUID, uuid1, uuid4
from datetime import timedelta as td
from sqlalchemy.engine import Engine

class SchedulerExecutionError(TypeError):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

class DuplicateUserError(BaseException):
    pass


def asynchronise(func, *args, **kwargs) -> Awaitable:
    """Handy function to convert a regular function to an async one."""
    function = partial(func, *args, **kwargs)
    async def af():
        return function()

    return af

class Job:
    """
    Props:
        is_due (bool):
            True if the task is currently due or overdue
        overdue_amount (datetime.timedelta):
            The difference between current time and the due time
        next_run (datetime):
            The next time this job is due to be run
    """
    KEYS = [
        "id",
        "user_id",
        "type",
        "interval",
        "last_run",
        "next_run"
    ]

    def __init__(
            self,
            user_id: UUID,
            interval: td,
            *args,
            name: str = None,
            job_id: UUID = None,
            job_type = None,
            **kwargs) -> None:
        """
        Interval will eventually be defined as fastapischeduler.Interval class
        Type will eventually be defined as fastapischeduler.Type class
        """
        self.id: UUID = uuid1()
        self.user_id: UUID = user_id
        self.type = job_type
        self.name = name

        self._interval: td = interval
        self._last_run: dt = dt.now()
        self._calculate_next_run()
        return None

    def __str__(self):
        return f"<Job: {self.name} interval: {self.interval} next due at: {self._next_run}{ ' (overdue by: ' + str(self.overdue_amount) + ')' if self.is_due else ''}>"

    def __getitem__(self, key):
        values = [
            self.id,
            self.user_id,
            self.type, 
            self._interval,
            self._last_run,
            self._next_run
        ]
        all_attrs = zip(self.KEYS, values)
        return dict(all_attrs)[key]

    def _calculate_next_run(self) -> dt:
        self._next_run = self._last_run + self.interval
        return self._next_run

    def keys(self):
        return self.KEYS

    @property
    def is_due(self):
        if self._next_run < dt.now():
            return True
        else:
            return False

    @property
    def overdue_amount(self):
        return dt.now() - self._next_run

    @property
    def next_run(self):
        return self._next_run

    # There is potential to add these methods in for the future.
    # @property
    # def last_result(self) -> Any:
    #     return self._last_result

    # @last_result.setter
    # def set_last_result(self, result):
    #     self._last_result = result

    def run(self) -> None:
        self._last_run = dt.now()
        self._calculate_next_run()
        return None

    @next_run.setter
    def next_run_setter(self, proposed_next_run: dt) -> None:
        if proposed_next_run < dt.now():
            raise ValueError(f"Cannot set a next run time in the past.")
        else:
            self._next_run = proposed_next_run
        return None

    @property
    def interval(self):
        return self._interval

    @interval.setter
    def interval_setter(self, proposed_interval: td):
        self._interval = proposed_interval


class Scheduler:
    """Scheduling class for use with FastAPI.

    It is based on a prototype function which all jobs are run in the
    context of.

    Attributes:
        is_running (bool):  Indicates the current working state.

    Methods:
        start:      Joins the scheduler into the event loop.
        stop:       Removes the scheduling task from the event loop.
        execute:    Executes a given job for the given scheduler.
    """
    STOPPED = 0
    RUNNING = 1
    PAUSED = 2
    NEXT_CALL_DELAY = 0.5

    def __init__(self,
            func: callable,
            db_engine: Engine,
            log: logging.Logger = None) -> None:
        """
        Args:
            func (callable):    Each task is called in the context of this fucntion, with
                                the Job instances properties passed as keyword arguments.
                                More documentation on constructing these functions will
                                be added.
            db_engine (Engine): SQLAlchemy engine used to execute database queries and
                                persist jobs across restarts.
            log (Logger):       Optional external log, defaults to module level.
        Returns:
            None
        """
        self._log: logging.Logger = log
        if not self._log:
            self._log = logging.getLogger(__name__)
            self._log.setLevel(logging.DEBUG)

        self._func = func

        self._jobs: List[Job] = self._fetch_from_database()
        self._loop: asyncio.BaseEventLoop = asyncio.get_event_loop()

    async def _execute_due(self) -> List[Any]:
        """Iterate over all the due jobs and (attempt to) capture errors along the way
        This is basically:

            async for job in self._fetch_due_jobs():
                self.execute(job)

        However gather means we can retain the results from each run of jobs executed,
        which might be important.

        Returns:
            result (List[Any]): A list of all the return values of the jobs executed in
                                this session.
        """
        result = await asyncio.gather(*[self.execute(job) async for job in self._fetch_due_jobs()])
        # Alert listeners to this "result"
        return result

    async def _task_loop(self) -> None:
        """ An asyncio task which calls necessary functions and listens for events. """
        while self._running:
            self._log.info("Checking for new jobs...")
            results = await self._execute_due()
            if results:
                self._log.info(results)
            else:
                self._log.debug(results)
            # Some magic here to see when the next nearest due job is and then sleep until then
            await asyncio.sleep(self.NEXT_CALL_DELAY)

    async def _fetch_due_jobs(self):
        """async generator to fetch only jobs which are due or overdue"""
        for job in self._jobs:
            if job.is_due:
                yield job

    @abstractmethod
    def _fetch_from_database(self) -> None:
        """
        Resets the self._jobs store with information from the server.
        """
        # try:
        #     self._jobs = self._db.fetch()
        # except Exception as e:
        #     self._log.error(e)
        job1 = Job(uuid4(), td(seconds=2))
        job2 = Job(uuid4(), td(seconds=5))
        job3 = Job(uuid4(), td(seconds=1))
        job4 = Job(uuid4(), td(seconds=10))
        return []

    @abstractmethod
    async def _sync_with_database(self) -> None:
        """
        Pushes the current state into the database and then fecthes it back.
        """
        # try:
        #     self._db.insert(self._jobs)
        #     self._fetch_from_database()
        # except Exception as e:
        #     self._log(e)
        # return

    @property
    def running(self):
        return self._running

    def execute(self, job: Job) -> asyncio.Task:
        """Executes a job object in the current scheduler

        Args:
            job (Job): A Job object to be executed
        Returns:
            task (asyncio.Task): Task as a futures object
        """
        self._log.debug(job)
        # Make this an optional step, and add in checking to see if the function
        # is async ready before doing it.
        task_func = asynchronise(self._func, **job)
        task = self._loop.create_task(task_func())

        job.run()
        return task

    def add(self, job: Job) -> UUID:
        """Adds a job to the jobs queue"""
        if job.user_id in [job["user_id"] for job in self._jobs]:
            raise DuplicateUserError("The given user is already present as a job.")
        self._jobs.append(job)
        self._loop.create_task(self._execute_due())
        return job.id

    def pop_user(self, user_id) -> Job:
        """Removes a job from the stack by user id"""
        for index, job in enumerate(self._jobs):
            if job["user_id"] == user_id:
                return self._jobs.pop(index)
        raise ValueError("The given user id does not exist.")

    def start(self) -> None:
        """
        The scheduler is running and can now execute jobs.
        """
        try:
            self._loop = asyncio.get_running_loop()
        except RuntimeError:
            self._loop = asyncio.get_event_loop()

        self._running = self.RUNNING
        self._loop.create_task(self._task_loop())

    def stop(self) -> None:
        """Set the state of the scheduler to stopped"""
        self._running = self.STOPPED
        return None

"""
base_wait_time = 1000ms
while true -> iterate(jobs_list) -> job_due : execute job :then: sleep(base_wait_time)
"""
def simplefunction(*args, **kwargs) -> int:
    from random import randint
    return str(f"user_id: {kwargs['user_id']}")
    # return randint(1, 999999999)

sch = Scheduler(simplefunction, None)