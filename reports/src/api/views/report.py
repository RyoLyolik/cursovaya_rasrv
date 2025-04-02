from typing import Annotated

from providers.report import get_report_service
from schemes.report import GenerateRequest
from services.report import ReportService

from fastapi import APIRouter, Depends


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


@router.get('/file')
async def download(
    report_service: Annotated[ReportService, Depends(get_report_service)]
):
    await report_service.get()
    return ...


@router.get('/file/list')
async def list_(
    report_service: Annotated[ReportService, Depends(get_report_service)]
):
    await report_service.list()
    return ...
