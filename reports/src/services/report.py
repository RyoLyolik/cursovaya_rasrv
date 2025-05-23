
import os

from datetime import timedelta
from uuid import uuid4

from repositories.data import DataRepository
from repositories.report import ReportRepository
from schemes.report import GenerateRequest

import matplotlib.pyplot as plt

from matplotlib.backends.backend_pdf import PdfPages


class ReportService:
    def __init__(self, report_repo: ReportRepository, data_repo: DataRepository):
        self.report_repo = report_repo
        self.data_repo = data_repo

    async def generate(self, req: GenerateRequest):
        tsp = [
            50, 60, 75, 100,  # 4
            125, 160, 200, 250,  # 8
            310, 400, 490, 570,  # 12
            630, 690, 740, 770,  # 16
            795, 815, 840, 855,  # 20
            870, 890, 910, 925,  # 24
            930, 940, 937, 930,  # 28
            915, 910, 890, 875,  # 32
            860, 845, 830, 800,  # 36
            760, 730, 650, 425,  # 40
            300, 225, 150, 120,  # 44
            80, 60,  # 46
        ]
        psp = [
            36, 37, 38, 39,  # 4
            39, 38, 36.8, 35,  # 8
            34, 31, 26, 24,  # 12
            22, 18, 14, 13,  # 16
            11, 10, 9, 8.7,  # 20
            8.5, 8.4, 8.3, 8.2,  # 24
            8.1, 8, 8, 8,  # 28
            8, 8, 8, 8,  # 32
            7.95, 7.9, 7.85, 7.7,  # 36
            7.4, 7, 6.6, 6.2,  # 40
            5.8, 5.8, 5.8, 5.8,  # 44
            5.8, 5.8,  # 46
        ]
        req.timefrom -= timedelta(hours=3)
        req.timeto -= timedelta(hours=3)
        temp_data = await self.fill_data(
            tablename='temperature',
            grouping=req.grouping,
            timefrom=req.timefrom,
            timeto=req.timeto,
            target=tsp,
            size=46
        )
        press_data = await self.fill_data(
            tablename='pressure',
            grouping=req.grouping,
            timefrom=req.timefrom,
            timeto=req.timeto,
            target=psp,
            size=46
        )
        flap_data = await self.fill_data(
            tablename='flap',
            grouping=req.grouping,
            timefrom=req.timefrom,
            timeto=req.timeto,
            target=[0] * 10,
            size=10
        )
        req.timefrom += timedelta(hours=3)
        req.timeto += timedelta(hours=3)
        ts = req.timefrom.strftime('%d-%m-%Y %H:%M:%S')
        te = req.timeto.strftime('%d-%m-%Y %H:%M:%S')
        object_name = f'отчет за {ts} — {te}.{str(uuid4())[:4]}.pdf'
        filename = object_name
        with PdfPages(filename) as pdf:
            # Первая страница с заголовком
            plt.figure(figsize=(8, 6))
            plt.axis('off')
            plt.text(0.5, 0.5, f'Отчет за\n{ts} — {te}', ha='center', va='center', fontsize=24)
            pdf.savefig()
            plt.close()

            plt.figure(figsize=(8, 6))
            plt.axis('off')
            plt.text(0.5, 0.5, 'Температурные параметры', ha='center', va='center', fontsize=24)
            pdf.savefig()
            plt.close()
            for pos in temp_data:
                tdata = temp_data[pos]
                dates = tdata[0]
                values = tdata[1]
                wanted = [tsp[pos - 1]] * 2
                plt.figure(figsize=(10, 6))
                plt.xlim((req.timefrom, req.timeto))
                plt.plot(dates, values, marker='o', linestyle='-', color='b', label='Записанные значения')
                plt.plot([req.timefrom, req.timeto], wanted, marker='o', linestyle='-', color='g', label='Желаемые значения')
                plt.title(f'Значения в позиции {pos}', fontsize=16)
                plt.xlabel('Время', fontsize=12)
                plt.ylabel('Значение параметра', fontsize=12)
                plt.grid(True)
                plt.xticks(rotation=45)
                plt.tight_layout()
                pdf.savefig()
                plt.close()

            plt.figure(figsize=(8, 6))
            plt.axis('off')
            plt.text(0.5, 0.5, 'Параметры давления', ha='center', va='center', fontsize=24)
            pdf.savefig()
            plt.close()
            for i, pos in enumerate(press_data.keys()):
                pdata = press_data[pos]
                dates = pdata[0]
                values = pdata[1]
                wanted = [psp[i]] * 2
                if not dates:
                    continue
                plt.figure(figsize=(10, 6))
                plt.xlim((req.timefrom, req.timeto))
                plt.plot(dates, values, marker='o', linestyle='-', color='b', label='Записанные значения')
                plt.plot([req.timefrom, req.timeto], wanted, marker='o', linestyle='-', color='g', label='Желаемые значения')
                plt.title(f'Значения в позиции {i + 1}', fontsize=16)
                plt.xlabel('Время', fontsize=12)
                plt.ylabel('Значение параметра', fontsize=12)
                plt.grid(True)
                plt.xticks(rotation=45)
                plt.tight_layout()
                pdf.savefig()
                plt.close()

            plt.figure(figsize=(8, 6))
            plt.axis('off')
            plt.text(0.5, 0.5, 'Положения заслонок', ha='center', va='center', fontsize=24)
            pdf.savefig()
            plt.close()
            for i, pos in enumerate(flap_data.keys()):
                fdata = flap_data[pos]
                dates = fdata[0]
                values = fdata[1]
                plt.figure(figsize=(10, 6))
                plt.xlim((req.timefrom, req.timeto))
                plt.ylim((0, 100))
                plt.plot(dates, values, marker='o', linestyle='-', color='b')
                plt.title(f'Значения в позиции {i + 1}', fontsize=16)
                plt.xlabel('Время', fontsize=12)
                plt.ylabel('Значение параметра', fontsize=12)
                plt.grid(True)
                plt.xticks(rotation=45)
                plt.tight_layout()
                pdf.savefig()
                plt.close()
        self.report_repo.fadd(filename, object_name)
        try:
            os.remove(filename)
        except Exception:
            print("failed to remove file")

    async def fill_data(self, tablename, grouping, timefrom, timeto, target, size=46):
        to_fill = {i: [[], [], target[i - 1]] for i in range(1, size + 1)}
        for position in to_fill:
            result = await self.data_repo.get_min_max_grouped(tablename, grouping, timefrom, timeto, position, target[position - 1])
            for record in result:
                to_fill[position][0].append(record.timestamp+timedelta(hours=3))
                to_fill[position][1].append(record.value)
        return to_fill

    def get(self, object_name):
        return self.report_repo.get(object_name)

    def list(self):
        res = self.report_repo.list()
        ans = []
        for obj in res:
            ans.append({
                'filename': obj.object_name,
                'size': obj.size,
                'creation_date': obj.last_modified,
            })
        return ans
