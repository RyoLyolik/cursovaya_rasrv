from fastapi import APIRouter


router = APIRouter(
    prefix='/report',
    tags=['auth'],
    responses={
        401: {},
        404: {},
    },
)
