
from schemes.data import Filter, FoundData

from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncSession


class DataRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_filter(self, filt: Filter, limit=True) -> list[FoundData]:
        query_start = f'''
        SELECT *
        FROM {filt.datatype}'''
        qc = 0
        query = ''
        if filt.timefrom and filt.timeto:
            query += f' (timestamp BETWEEN \'{filt.timefrom.isoformat()}\' AND \'{filt.timeto.isoformat()}\')'
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
            final += ' WHERE' + query + ' ORDER BY timestamp' + ' LIMIT 500' * limit + ';'  # TODO make pagination after
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

    async def get_min_max_grouped(self, tablename, grouping, timefrom, timeto, position, target) -> list[FoundData]:
        query = f'''
WITH hourly_peaks AS (
    SELECT
        DATE_TRUNC('{grouping}', timestamp) AS timeunit,
        position,
        MAX(ABS(value-{target})) AS peak_value
    FROM
        {tablename}
    WHERE
        timestamp BETWEEN '{timefrom}' AND '{timeto}'
        AND position = {position}
    GROUP BY
        DATE_TRUNC('{grouping}', timestamp),
        position
)
SELECT
    hp.timeunit,
    hp.position,
    hp.peak_value,
    t.timestamp AS peak_time
FROM
    hourly_peaks hp
JOIN
    temperature t
    ON DATE_TRUNC('{grouping}', t.timestamp) = hp.timeunit
    AND t.position = hp.position
    AND t.value = hp.peak_value
ORDER BY
    hp.timeunit,
    hp.position;
'''
        result = await self.session.execute(text(query))
        rows = result.all()
        result = []
        for row in rows:
            fdata = FoundData.model_construct(
                timestamp=row[3],
                position=row[1],
                value=row[2],
            )
            result.append(fdata)
        return result
