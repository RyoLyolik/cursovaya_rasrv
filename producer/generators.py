from datetime import datetime, timezone
from typing import Generator
from value_generator import value_generator


ANOMALY_CHANGE = 0.000001


def temperature_generator() -> Generator[dict, None, None]:
    set_points = [
        50, 60, 75, 100, # 4
        125, 160, 200, 250, # 8
        310, 400, 490, 570, # 12
        630, 690, 740, 770, # 16
        795, 815, 840, 855, # 20
        870, 890, 910, 925, # 24
        930, 940, 937, 930, # 28
        915, 910, 890, 875, # 32
        860, 845, 830, 800, # 36
        760, 730, 650, 425, # 40
        300, 225, 150, 120, # 44
        80, 60 # 46
    ]
    generator = setup_generator(set_points)
    for position, record in generator:
        yield process_data(record, 'temperature', position)


def pressure_generator() -> Generator[dict, None, None]:
    set_points = [
        36, 37, 38, 39, # 4
        39, 38, 36.8, 35, # 8
        34, 31, 26, 24, # 12
        22, 18, 14, 13, # 16
        11, 10, 9, 8.7, # 20
        8.5, 8.4, 8.3, 8.2, # 24
        8.1, 8, 8, 8, # 28
        8, 8, 8, 8, # 32
        7.95, 7.9, 7.85, 7.7, # 36
        7.4, 7, 6.6, 6.2, # 40
        5.8, 5.8, 5.8, 5.8, # 44
        5.8, 5.8, # 46
    ]
    generator = setup_generator(set_points)
    for position, record in generator:
        yield process_data(record, 'pressure', position)


def flap_generator() -> Generator[dict, None, None]:
    set_points = [50]*10
    generator = setup_generator(set_points)
    for position, record in generator:
        yield process_data(record, 'flap', position)


def setup_generator(setpoints: list[float]) -> Generator[tuple[int, int], None, None]:
    generators = [
        value_generator(
            set_point=sp,
            noise_level=sp/100,
            anomaly_chance=ANOMALY_CHANGE,
            anomaly_size=sp,
        )
        for sp in setpoints
    ]
    while True:
        for i, generator in enumerate(generators):
            yield i+1, next(generator)


def process_data(value: float, type: str, position: int) -> dict:
    tmstmp = datetime.now(tz=timezone.utc)
    return {
        'timestamp': tmstmp.isoformat(),
        'record_id': tmstmp.timestamp(),
        'record_type': type,
        'value': value,
        'pos_id': position
    }
