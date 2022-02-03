from datetime import datetime as dt
from datetime import timedelta as td
from typing import Any, Dict, List
from uuid import UUID, uuid1


class Job:
    """Job class.

    Takes args of the user which it belongs to and a time interval to run at.

    Args:
        user_id (UUID):
            The id of the user this job belongs to
        interval (timedelta):
            The interval at which this job should run

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
        "job_type",
        "interval",
        "last_run",
        "next_run",
    ]

    def __init__(
            self,
            user_id: UUID,
            interval: td,
            *args,
            name: str = None,
            job_id: UUID = None,
            job_type = None,
            mail = None,
            **kwargs) -> None:
        """
        Interval will eventually be defined as fastapischeduler.Interval class
        Type will eventually be defined as fastapischeduler.Type class
        """
        self.id: UUID = uuid1()
        self.user_id: UUID = user_id
        self.job_type = job_type
        self.name = name
        self.mail = mail

        self._interval: td = interval
        self._last_run: dt = dt.now()
        self._calculate_next_run()
        return None

    def __str__(self):
        return f"<Job: {self.name} interval: {self.interval} next due at: {self._next_run}{ ' (overdue by: ' + str(self.overdue_amount) + ')' if self.is_due else ''}>"

    def __getitem__(self, key: str) -> Dict[str, Any]:
        values = [
            self.id,
            self.user_id,
            self.job_type,
            self._interval,
            self._last_run,
            self._next_run,
        ]
        all_attrs = zip(self.KEYS, values)
        return dict(all_attrs)[key]

    def _calculate_next_run(self) -> dt:
        self._next_run = self._last_run + self.interval
        return self._next_run

    def keys(self) -> List[str]:
        return self.KEYS

    def run(self) -> None:
        """Called by a scheduler to mark the job as run and update the next due time"""
        self._last_run = dt.now()
        self._calculate_next_run()
        return None

    @property
    def is_due(self) -> bool:
        """Returns True if the job is due or overdue"""
        return self._next_run < dt.now()

    @property
    def overdue_amount(self) -> td:
        """Returns timedelta object of ```now - time of next run```"""
        return dt.now() - self._next_run

    @property
    def next_run(self) -> dt:
        """The time when the job next needs to be run"""
        if not self._next_run:
            self._calculate_next_run()
        return self._next_run

    @next_run.setter
    def next_run_setter(self, proposed_next_run: dt) -> None:
        current_time = dt.now()
        if proposed_next_run < current_time:
            raise ValueError(f"Cannot set a next run time in the past.")
        else:
            self._next_run = proposed_next_run
        return None

    @property
    def interval(self) -> td:
        """The interval at which this job should be run"""
        return self._interval

    @interval.setter
    def interval_setter(self, proposed_interval: td):
        self._interval = proposed_interval

    """There is potential to add these methods in for the future.
    @property
    def last_result(self) -> Any:
        return self._last_result

    @last_result.setter
    def set_last_result(self, result):
        self._last_result = result"""
