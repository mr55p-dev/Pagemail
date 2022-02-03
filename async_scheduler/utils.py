from functools import partial
from typing import Awaitable

"""Utility functions which are used by the scueduler"""

def partial_asynchronise(func, *args, **kwargs) -> Awaitable:
    """Handy function to convert a regular function into a partial ```async``` one."""
    function = partial(func, *args, **kwargs)
    async def af():
        return function()

    return af

def asynchronise(func):
    async def af(*args, **kwargs):
        return func(*args, **kwargs)
    return af

def relu(value: float) -> float:
    if value >= 0.0:
        return float(value)
    else:
        return float(0.0)
