FROM python:3.12 as python

ENV PYTHONOPTIMIZE=1 \
    PYTHONPATH=/usr/src \
    PYTHONUNBUFFERED=1

WORKDIR ${PYTHONPATH}

RUN python3 -m venv .venv
RUN . ${PYTHONPATH}/.venv/bin/activate

COPY ./requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY . .
CMD ["python", "producer.py"]
