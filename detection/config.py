from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    REDIS_HOST: str = "localhost"
    REDIS_PORT: int = 6379
    REDIS_DB: int = 0

    DB_HOST: str = "localhost"
    DB_PORT: int = 5432
    DB_NAME: str = "adr"
    DB_USER: str = "adr"
    DB_PASSWORD: str = "adr"

    ML_RETRAIN_INTERVAL_SEC: int = 300
    ML_WINDOW_SIZE: int = 1000

    class Config:
        env_file = ".env"


settings = Settings()
