from repositories.report import ReportRepository


class ReportService:
    def __init__(self, repo: ReportRepository):
        self.repo = repo