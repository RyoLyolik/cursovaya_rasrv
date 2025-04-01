from datetime import datetime
from typing import Literal

from pydantic import BaseModel


class FoundData(BaseModel):
    position: int
    value: float
    timestamp: datetime


class Filter(BaseModel):
    datatype: Literal['temperature', 'pressure', 'flap']
    lower_equal: float | None
    greater_equal: float | None
    position: int | None
    timefrom: datetime | None
    timeto: datetime | None
