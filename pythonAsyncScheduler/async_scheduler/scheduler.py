import asyncio
import logging
import inspect
from abc import abstractmethod
from datetime import datetime as dt
from datetime import timedelta as td
from typing import Any, List
from uuid import UUID

# from sqlalchemy.engine import Engine

from async_scheduler.exceptions import (DuplicateUserError,
                                        PrototypeFunctionError, SchedulerRunningError)
from async_scheduler.job import Job
from async_scheduler.utils import asynchronise, relu


class Scheduler:
    """Scheduling class for use with FastAPI.

    It is based on a prototype function which all jobs are run in the
    context of.

    Properties:
        is_running (bool):
            Indicates the current working state.
        function (callable):
            Returns the current context function.

    Methods:
        start:
            Joins the scheduler into the event loop.
        stop:
            Removes the scheduling task from the event loop.
        execute:
            Executes a given job for the given scheduler.
        add (Job):
            Add a job to the scheduler.
        pop_user (Job):
            Remove a job from the scheduler by its ```user_id``` and return it.

    ```func (callable)``` is required to be a callable which can take each of
    the ```[user_id, job_type]``` of a ```Job``` as ```kwargs```.
    """
    STOPPED = 0
    RUNNING = 1
    PAUSED = 2
    next_call_delay: float = 10.0
    BASE_CALL_DELAY: float = 10.0

    def __init__(self,
            # db_engine: Engine,
            target_func: callable = None,
            log: logging.Logger = None) -> None:
        """
        Args:
            func (callable):    Each task is called in the context of this fucntion.
            db_engine (Engine): SQLAlchemy engine used to execute database queries and
                                persist jobs across restarts.
            log (Logger):       Optional external log, defaults to module level.
        Returns:
            None
        """
        if log:
            self._log: logging.Logger = log
        else:
            self._log = logging.getLogger(__name__)
            self._log.setLevel(logging.DEBUG)

        self._func = None
        self._running = self.STOPPED
        self._jobs: List[Job] = self._fetch_from_database()
        self._loop: asyncio.BaseEventLoop = asyncio.get_event_loop()
        self._lock: asyncio.Lock = asyncio.Lock()

        self.function = target_func

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

            await self._calculate_sleep()
            self._log.debug(f"Sleeping for: {self.next_call_delay} seconds")

            await asyncio.sleep(self.next_call_delay)

    async def _fetch_due_jobs(self):
        """async generator to fetch only jobs which are due or overdue"""
        for job in self._jobs:
            if job.is_due:
                yield job

    async def _calculate_sleep(self) -> float:
        """Lookup how long to wait until running the next jobs"""
        if self._jobs:
            next_job: Job = sorted(self._jobs, key=lambda i: i["next_run"])[0]
            due_time: dt = next_job.next_run
            sleep_time: td = due_time - dt.now()
            self.next_call_delay = relu(sleep_time.total_seconds())
        else:
            self.next_call_delay = self.BASE_CALL_DELAY

        return self.next_call_delay

    async def _real_set_func(self, newfunc: callable) -> None:
        async with self._lock:
            if inspect.iscoroutinefunction(newfunc):
                self._func = newfunc
            else:
                self._func = asynchronise(newfunc)

    @abstractmethod
    def _fetch_from_database(self) -> List[Job]:
        """Resets the self._jobs store with information from the server."""
        # try:
        #     self._jobs = self._db.fetch()
        # except Exception as e:
        #     self._log.error(e)
        return []

    @abstractmethod
    async def _sync_with_database(self) -> None:
        """Pushes the current state into the database and then fecthes it back."""
        # try:
        #     self._db.insert(self._jobs)
        #     self._fetch_from_database()
        # except Exception as e:
        #     self._log(e)
        # return

    @property
    def running(self) -> bool:
        """Returns ```True``` if running and ```False``` if stopped"""
        return self._running == self.RUNNING

    @property
    def function(self) -> callable:
        """Returns the function which is used by the scheduler"""
        if self._func:
            return self._func
        else:
            return None
        # return self._func or None

    @function.setter
    def function(self, newfunc: callable) -> None:
        if self.running:
            raise SchedulerRunningError(
                "Cannot reassign target function while the scheduler is running.")
        else:
            self._loop.create_task(self._real_set_func(newfunc))

    def execute(self, job: Job) -> asyncio.Task:
        """Executes a job object in the current scheduler

        Args:
            job (Job): A Job object to be executed
        Returns:
            task (asyncio.Task): Task as a futures object
        """
        self._log.debug(job)

        # Extract the arguments from the job
        # Make these arguments tunable via a flag of some kind.
        job_args = {i: job[i] for i in ["user_id", "job_type"] if i in job.keys()}

        # Partialise and asynchronise the function
        # Make this an optional step, and add in checking to see if the function
        # is async ready before doing it.
        # task_func = partial_asynchronise(self._func, **job_args)

        # Add the function as a task
        task = self._loop.create_task(self._func(**job_args))

        job.run()
        return task

    def add(self, job: Job) -> UUID:
        """Adds a job to the jobs queue

        Args:
            job (Job): Take a job and insert it into the scheduler.

        Returns:
            job_id (UUID): The assigned job_id.
        """
        if job.user_id in [job["user_id"] for job in self._jobs]:
            raise DuplicateUserError("The given user is already present as a job.")
        self._jobs.append(job)
        self._loop.create_task(self._execute_due())
        return job.id

    def pop_user(self, user_id) -> Job:
        """Removes a job from the stack by user id

        Args:
            user_id (UUID): User id of the job to pop

        Returns:
            job (Job): The full job which has been popped
        """
        for index, job in enumerate(self._jobs):
            if job["user_id"] == user_id:
                return self._jobs.pop(index)
        raise ValueError("The given user id does not exist.")

    def start(self) -> None:
        """The scheduler is running and can now execute jobs."""
        if not self._func:
            raise PrototypeFunctionError(
                "Cannot start the scheduler without a contextual function.")

        self._running = self.RUNNING
        self._loop.create_task(self._task_loop())

    def stop(self) -> None:
        """Set the state of the scheduler to stopped"""
        self._running = self.STOPPED
