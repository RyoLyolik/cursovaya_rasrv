import os
import re

from typing import Literal

import yaml

from pydantic import BaseModel
from pydantic_settings import BaseSettings


class AppConfig(BaseModel):
    environment: Literal['test', 'dev', 'prod']
    host: str
    port: int
    workers: int
    timeout: int
    hostname: str


class LoggingConfig(BaseModel):
    log_level: str
    log_format: str


class DataStorages(BaseModel):
    main: list[str]


class Settings(BaseSettings):
    app: AppConfig
    logging: LoggingConfig
    data_storages: DataStorages


def substitute_env_variables(value: str) -> str:
    matches = re.findall(r'\${(.*?)}', value)

    for match in matches:
        env_value = os.getenv(match)
        if env_value is not None:
            value = value.replace(f'${{{match}}}', env_value)
        elif os.getenv('ENVIRONMENT') != 'test':
            raise Exception(f'Error: Environment variable {match} not found.')

    return value


def load_yaml_settings(filepath: str) -> Settings:
    if filepath is None:
        filepath = 'local.yaml'
    with open(filepath, 'r') as file:
        data = yaml.safe_load(file)

    def recursive_substitute(d):
        if isinstance(d, dict):
            return {key: recursive_substitute(value) for key, value in d.items()}
        elif isinstance(d, list):
            return [recursive_substitute(item) for item in d]
        elif isinstance(d, str):
            return substitute_env_variables(d)
        else:
            return d

    data = recursive_substitute(data)

    return Settings(**data)


settings = load_yaml_settings(os.getenv('SETTINGS_FILE'))
