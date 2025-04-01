from providers.db import get_async_session
from repositories.report import ReportRepository
from services.report import ReportService

from minio import Minio


def get_report_storage():
    return Minio(
        host='minio',
        port='9000',
        access_key='minioadmin',
        secret_key='minioadmin',
        main_bucket="files",
        secure=False,
    )


def get_report_service():
    client = get_report_storage()
    session = get_async_session()
    repo = ReportRepository(client, session)
    service = ReportService(repo)
    return service
