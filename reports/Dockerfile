FROM python:3.12 as python

# Задаём переменные окружения.
# UV_COMPILE_BYTECODE — включает компиляцию файлов Python в байт-код после установки.
# UV_LINK_MODE — меняет способ установки пакетов из глобального кэша.
#   Вместо создания жёстких ссылок, файлы пакета копируются в директорию  виртуального окружения `site-packages`.
#   Это необходимо для будущего копирования изолированной `/app` директории из  стадии `build` в финальный Docker-образ.
# UV_PYTHON — фиксирует версию Python.
# UV_PROJECT_ENVIRONMENT — указывает путь к виртуальному окружению Python.
# UV_PYTHON_DOWNLOADS — отключает автоматическую загрузку отсутствующих версий Python.
# https://docs.astral.sh/uv/configuration/environment/
# PYTHONOPTIMIZE — убирает инструкции `assert` и код, зависящий от значения  константы `__debug__`, при компиляции файлов Python в байт-код.
# PYTHONUNBUFFERED Эта переменная отключает буферизацию стандартных потоков вывода (stdout) и ошибок (stderr).
# PYTHONPATH - Эта переменная задает дополнительные пути, по которым Python будет искать модули и пакеты.
# https://docs.python.org/3/using/cmdline.html#environment-variables
ENV PYTHONOPTIMIZE=1 \
    PYTHONPATH=/usr/prj/src \
    PYTHONUNBUFFERED=1 \
    UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_PYTHON="python3.12" \
    UV_PROJECT_ENVIRONMENT=${PYTHONPATH}/.venv \
    UV_PYTHON_DOWNLOADS=never

WORKDIR ${PYTHONPATH}/../

RUN python3 -m venv .venv
RUN . ${PYTHONPATH}/../.venv/bin/activate

# Копируем утилиту `uv` из официального Docker-образа.
# https://github.com/astral-sh/uv/pkgs/container/uv
# опция `--link` позволяет переиспользовать слой, даже если предыдущие слои изменились.
# https://docs.docker.com/reference/dockerfile/#copy---link
COPY --link --from=ghcr.io/astral-sh/uv:0.4 /uv /usr/local/bin/uv

COPY ./pyproject.toml pyproject.toml
RUN uv sync

COPY ./local.yaml local.yaml
COPY ./src src
