from datetime import datetime
from typing import Annotated

from providers.data_service import get_data_service
from schemes.data import Filter, FoundData
from services.data import DataService

from fastapi import APIRouter, Depends


router = APIRouter(
    prefix='/data',
    tags=['auth'],
    responses={
        401: {},
        404: {},
    },
)


@router.get('/find', response_model=list[FoundData])
async def search_temperature(
    data_service: Annotated[DataService, Depends(get_data_service)],
    datatype: str,
    le: float | None = None,
    ge: float | None = None,
    position: int | None = None,
    timefrom: datetime | None = None,
    timeto: datetime | None = None,
):
    filt = Filter.model_construct(
        datatype=datatype,
        lower_equal=le,
        greater_equal=ge,
        position=position,
        timefrom=timefrom,
        timeto=timeto,
    )
    return await data_service.find_data(filt)
