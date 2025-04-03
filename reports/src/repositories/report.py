

from minio import Minio
from sqlalchemy.ext.asyncio import AsyncSession


class ReportRepository:
    def __init__(self, client: Minio, session: AsyncSession):
        self.client = client
        self.bucket = 'reports'
        if not self.client.bucket_exists(self.bucket):
            self.client.make_bucket(self.bucket)

    def fadd(self, fp, object_name):
        self.client.fput_object(
            bucket_name=self.bucket,
            object_name=object_name,
            file_path=fp
        )

    def get(self, object_name):
        return self.client.get_object(self.bucket, object_name).data

    def list(self):
        return self.client.list_objects(self.bucket)
