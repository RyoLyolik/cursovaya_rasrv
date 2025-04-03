import requests
import time
from concurrent.futures import ThreadPoolExecutor
import statistics

def make_request(url):
    start_time = time.time()
    try:
        response = requests.get(url)
        response.raise_for_status()
        elapsed_time = (time.time() - start_time) * 1000
        return elapsed_time
    except requests.exceptions.RequestException as e:
        print(f"Ошибка при запросе к {url}: {e}")
        return None

def run_load_test(url, num_requests, num_threads=10):
    times = []
    
    with ThreadPoolExecutor(max_workers=num_threads) as executor:
        futures = [executor.submit(make_request, url) for _ in range(num_requests)]
        
        for future in futures:
            result = future.result()
            if result is not None:
                times.append(result)

    print(f"Запросов выполнено: {len(times)}/{num_requests}")
    print(f"Среднее время ответа: {statistics.mean(times):.2f} мс")
    print(f"Медианное время ответа: {statistics.median(times):.2f} мс")
    print(f"Минимальное время: {min(times):.2f} мс")
    print(f"Максимальное время: {max(times):.2f} мс")
    print(f"Стандартное отклонение: {statistics.stdev(times):.2f} мс")

if __name__ == "__main__":
    BACKEND_URL = "http://localhost:7777/reports/list"
    NUM_REQUESTS = 100
    NUM_THREADS = 10
    print(f"Запуск нагрузочного теста для {BACKEND_URL}...")
    run_load_test(BACKEND_URL, NUM_REQUESTS, NUM_THREADS)
