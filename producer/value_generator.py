from random import random
from typing import Generator

def sign(x: float) -> int:
    if x < 0:
        return -1
    return 1

def value_generator(
    set_point: int,
    noise_level: float,
    anomaly_chance: float,
    anomaly_size: int
) -> Generator[int, None, None]:
    while True:
        if set_point is None:
            yield None
            continue

        action = random()
        noise = 2 * noise_level * (random()-0.5)
        anomaly = 0
        if action < anomaly_chance:
            anomaly = anomaly_size * 0.5 * sign(random()-0.5)
        value = set_point + anomaly + noise
        yield value
