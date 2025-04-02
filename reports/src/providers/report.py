from typing import Annotated
from fastapi import Depends
from providers.db import get_async_session
from repositories.data import DataRepository
from repositories.report import ReportRepository
from services.report import ReportService

from minio import Minio
from sqlalchemy.ext.asyncio import AsyncSession


def get_report_storage():
    return Minio(
        endpoint='minio:9000',
        access_key='minioadmin',
        secret_key='minioadmin',
        secure=False,
    )


def get_report_service(
    session: Annotated[AsyncSession, Depends(get_async_session)],
):
    client = get_report_storage()
    report_repo = ReportRepository(client, session)
    data_repo = DataRepository(session)
    service = ReportService(report_repo, data_repo)
    return service
