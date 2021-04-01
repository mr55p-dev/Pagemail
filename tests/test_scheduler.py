from API.helpers.scheduling import Scheduler, Job
from uuid import uuid1
from datetime import timedelta as td


job = Job(uuid1(), td(seconds=5))
job_dict_coerce = dict(job)
job_dict_compre = {**job}
job_dict_cast = {job}

assert job_dict_coerce == job_dict_cast == job_dict_compre
