
from minio import Minio
from sqlalchemy.ext.asyncio import AsyncSession


class ReportRepository:
    def __init__(self, client: Minio, session: AsyncSession):
        self.client = client
        self.bucket = 'reports'
