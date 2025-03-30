import asyncio
import json
import websockets
import logging
from asyncio import sleep as asleep
from generators import temperature_generator, pressure_generator, flap_generator
from sensor import Sensor


logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),  # Вывод в консоль
        logging.FileHandler('websocket_server.log')  # Запись в файл
    ]
)
logger = logging.getLogger(__name__)


sensors = [
    Sensor(temperature_generator()),
    Sensor(pressure_generator()),
    Sensor(flap_generator())
]
for sensor in sensors:
    sensor.start()


async def echo(websocket: websockets.serve):
    client_ip = websocket.remote_address[0]
    logger.info(f"Новое подключение от {client_ip}")
    try:
        while True:
            for sensor in sensors:
                await websocket.send(json.dumps(sensor.data))
            await asleep(0.05)

    except websockets.exceptions.ConnectionClosed:
        logger.info(f"Соединение с {client_ip} закрыто")
    except Exception as e:
        logger.error(f"Ошибка с клиентом {client_ip}: {str(e)}")


async def main():
    logger.info("Запуск WebSocket сервера на ws://localhost:2113")
    async with websockets.serve(echo, "0.0.0.0", 2113, server_header={
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type"
    }):
        logger.info("Сервер успешно запущен")
        await asyncio.Future()  # Бесконечное выполнение


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Сервер остановлен по запросу пользователя")
    except Exception as e:
        logger.error(f"Критическая ошибка сервера: {str(e)}")