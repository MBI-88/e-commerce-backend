from pydantic import BaseSettings


class BaseEnv(BaseSettings):
    ORIGINS: str
    CREDENTIALS: bool
    DEBUG: bool
    PORT_APP: int
    HOST_APP: str
    METHODS: str
    ALLOW_HOSTS: str
    DSN: str
    HASHKEY: str

    class Config:
        env_file = '.env'