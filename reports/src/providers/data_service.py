from typing import Annotated

from providers.db import get_async_session
from repositories.data import DataRepository
from services.data import DataService

from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession


async def get_data_service(
    session: Annotated[AsyncSession, Depends(get_async_session)],
):
    data_repo = DataRepository(session=session)
    return DataService(repo=data_repo)
