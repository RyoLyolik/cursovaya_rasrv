import os
import tempfile

from typing import Annotated

from providers.report import get_report_service
from schemes.report import GenerateRequest
from services.report import ReportService

from fastapi import APIRouter, BackgroundTasks, Depends
from fastapi.responses import FileResponse


router = APIRouter(
    prefix='/report',
    tags=['auth'],
    responses={
        401: {},
        404: {},
    },
)


@router.post('/generate')
async def upload(
    report_service: Annotated[ReportService, Depends(get_report_service)],
    generate_req: GenerateRequest
):
    await report_service.generate(generate_req)
    return {'status': 'OK'}


def delete_file(file_path: str):
    os.unlink(file_path)


@router.get('/file')
async def download(
    report_service: Annotated[ReportService, Depends(get_report_service)],
    filename: str,
    background_tasks: BackgroundTasks,
):
    tmp_path = ''
    iodata = report_service.get(object_name=filename)
    with tempfile.NamedTemporaryFile(delete=False) as tmp:
        tmp.write(iodata)
        tmp_path = tmp.name
    response = FileResponse(
        tmp_path, 200, media_type="application/octet-stream",
        filename=filename
    )
    background_tasks.add_task(delete_file, tmp_path)
    return response


@router.get('/file/list')
async def list_(
    report_service: Annotated[ReportService, Depends(get_report_service)]
):
    return report_service.list()
