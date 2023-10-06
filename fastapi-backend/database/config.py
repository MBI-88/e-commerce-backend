from sqlmodel import create_engine
from sqlmodel.ext.asyncio.session import AsyncSession


# Connection

database_connection_string = ""
connect_args = {"check_same_thread":False}
engine_url = ""

def setDatabaseConn(dsn:str) -> None:
    global database_connection_string 
    global engine_url
    database_connection_string = dsn
    engine_url = create_engine(database_connection_string,echo=True)

async def get_session() -> None:
    async with AsyncSession(engine_url) as session:
        yield await session