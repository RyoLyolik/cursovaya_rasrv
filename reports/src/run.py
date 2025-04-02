import logging

from contextlib import asynccontextmanager

import uvicorn

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import ORJSONResponse

from providers.db import get_async_engine


@asynccontextmanager
async def lifespan(app: FastAPI):
    # from src.utils.db import create_db_and_tables, delete_tables

    engine = get_async_engine()

    # await delete_tables(engine)
    # await create_db_and_tables(engine)

    yield

    await engine.dispose()


def create_app() -> FastAPI:
    from api.router import router

    app = FastAPI(
        default_response_class=ORJSONResponse,
        lifespan=lifespan,
    )

    app.add_middleware(
        CORSMiddleware,
        allow_origins=['*'],
        allow_credentials=True,
        allow_methods=['*'],
        allow_headers=['*'],
    )

    app.include_router(
        router,
    )

    return app


app = create_app()

if __name__ == '__main__':
    from core.config import settings

    logging.basicConfig(format=settings.logging.log_format)
    uvicorn.run(
        'run:app',
        host=settings.app.host,
        port=settings.app.port,
        workers=settings.app.workers,
        reload=True,
    )
