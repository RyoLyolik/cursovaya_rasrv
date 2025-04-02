from datetime import datetime
from typing import Literal

from pydantic import BaseModel


class GenerateRequest(BaseModel):
    timefrom: datetime
    timeto: datetime
    grouping: Literal['day', 'hour', 'minute', 'second']
