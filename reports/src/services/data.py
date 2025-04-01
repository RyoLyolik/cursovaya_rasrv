from repositories.data import DataRepository
from schemes.data import Filter, FoundData


class DataService:
    def __init__(self, repo: DataRepository):
        self.repo = repo

    async def find_data(self, filter: Filter) -> list[FoundData]:
        return await self.repo.get_by_filter(filter)
