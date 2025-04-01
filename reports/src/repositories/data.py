from datetime import datetime

from schemes.data import Filter, FoundData

from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncSession


class DataRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_filter(self, filt: Filter) -> list[FoundData]:
        query_start = f'''
        SELECT *
        FROM {filt.datatype}'''
        qc = 0
        query = ''
        if filt.timefrom and filt.timeto:
            query += f' (timestamp BETWEEN {filt.timefrom.isoformat()} AND {filt.timeto.isoformat()})'
            qc += 1
        elif filt.timefrom:
            query += f' timestamp >= \'{filt.timefrom.isoformat()}\''
            qc += 1
        elif filt.timeto:
            query += f' timestamp <= \'{filt.timeto.isoformat()}\''
            qc += 1

        if filt.lower_equal or filt.greater_equal:
            qc2 = 0
            if qc:
                query += ' AND'
            if filt.lower_equal:
                qc2 += 1
                query += f' value <= {filt.lower_equal}'
            if filt.greater_equal:
                if qc2:
                    query += ' AND'
                query += f' value >= {filt.greater_equal}'
            qc += 1
        if filt.position:
            if qc:
                query += ' AND'
            query += f' position = {filt.position}'
        final = query_start
        if query:
            final += ' WHERE' + query + ' LIMIT 500;'  # TODO make pagination after
        result = await self.session.execute(text(final))
        rows = result.all()
        result = []
        for row in rows:
            fdata = FoundData.model_construct(
                timestamp=row[1],
                position=row[2],
                value=row[3],
            )
            result.append(fdata)
        return result
