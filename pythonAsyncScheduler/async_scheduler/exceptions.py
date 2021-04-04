class SchedulerRunningError(BaseException):
    pass


class SchedulerExecutionError(TypeError):
    pass


class DuplicateUserError(BaseException):
    pass


class PrototypeFunctionError(BaseException):
    pass