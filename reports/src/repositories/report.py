
from typing import BinaryIO

from minio import Minio
from sqlalchemy.ext.asyncio import AsyncSession


class ReportRepository:
    def __init__(self, client: Minio, session: AsyncSession):
        self.client = client
        self.bucket = 'reports'

    def add(self, obj, data: BinaryIO):
        self.client.put_object(self.bucket, )

    def get(self, id):
        ...

    def list(self):
        ...
