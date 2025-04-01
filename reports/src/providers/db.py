from typing import AsyncGenerator

from core.config import settings

from sqlalchemy.ext.asyncio import (
    AsyncEngine,
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)


__engine = None


async def get_async_session() -> AsyncGenerator[AsyncSession, None]:
    async_session_maker = async_sessionmaker(get_async_engine(), expire_on_commit=False, class_=AsyncSession)
    async with async_session_maker.begin() as session:
        yield session


def get_async_engine(echo=False) -> AsyncEngine:
    global __engine
    if __engine is None:
        __engine = create_async_engine(
            settings.data_storages.main[0],
            echo=echo,
        )
    return __engine
