import json
import random
import time
from datetime import datetime, timezone
from confluent_kafka import Producer

# Конфигурация Kafka
KAFKA_BROKER = 'localhost:9092'
TOPIC_NAME = 'temperature-readings'

# Список локаций
LOCATIONS = ['room1', 'room2', 'room3']

def delivery_report(err, msg):
    """Callback-функция для обработки результатов отправки сообщения"""
    if err is not None:
        print(f'Message delivery failed: {err}')
    else:
        print(f'Message delivered to {msg.topic()} [{msg.partition()}]')

def generate_temperature_reading():
    """Генерация случайного показания температуры"""
    return {
        'location': random.choice(LOCATIONS),
        'temperature': round(random.uniform(18.0, 25.0), 1),
        'timestamp': datetime.now(timezone.utc).isoformat()
    }

def main():
    # Создаем продюсера
    conf = {
        'bootstrap.servers': KAFKA_BROKER,
        'client.id': 'temperature-producer'
    }
    producer = Producer(conf)

    try:
        while True:
            # Генерируем и отправляем данные каждые 3 секунды
            reading = generate_temperature_reading()
            producer.produce(
                TOPIC_NAME,
                key=reading['location'],
                value=json.dumps(reading),
                callback=delivery_report
            )
            
            # Ожидаем отправки сообщений
            producer.flush()
            
            print(f'Sent: {reading}')
            time.sleep(3)
            
    except KeyboardInterrupt:
        print("\nProducer stopped")
    finally:
        producer.flush()

if __name__ == '__main__':
    main()