from threading import Thread
from time import sleep
from typing import Generator


class Sensor:
    def __init__(self, generator: Generator[dict, None, None]):
        self.generator = generator
        self.data = None
    
    def load_data(self):
        self.data = next(self.generator)

    def update(self):
        while True:
            self.load_data()
            sleep(0.003)
    
    def start(self):
        th = Thread(target=self.update, args=(), daemon=True)
        th.start()
