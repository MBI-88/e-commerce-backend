from argparse import ArgumentParser, RawDescriptionHelpFormatter
import textwrap
import subprocess
import uvicorn as server
from fastapi import FastAPI
from services.bot.router import router as botRotuer
from database.config import setDatabaseConn
from middlewares.jwt import SECRET_KEY


# Application

app = FastAPI(description="E-commerce",title="D'Todo")

app.include_router(router=botRotuer, tags=["Bot"])



if __name__ == "__main__":
    command = ArgumentParser(
        description=textwrap.dedent(
            """
    -------------------------------------------------------------
    D'Todo fastapi backend
    -------------------------------------------------------------
    """
        ), prefix_chars='--', formatter_class=RawDescriptionHelpFormatter
    )

    command.add_argument('argument', type=str, nargs='?',
                         help="""run test shell database""")
    command.add_argument('-s', '--settings', type=str,
                         help="Import local settings to run application")
    command.add_argument('-P', '--port', type=int,
                         help="Application port")
    command.add_argument('-H', '--host', type=str,
                         help="Application host")
    command.add_argument('-f', '--file', type=str,
                         help="File to open in shell")
    selection = command.parse_args()

    match (selection.argument):
        case ('runserver'):
            # Run server
            host_app = ""
            port = 0
            host = ""
            dsn = ""
            if selection.settings == 'prod':
                from settings.pro import (
                    TrustedHostMiddleware, GZipMiddleware,
                    CORSMiddleware, HEADERS, ORIGINS_PRO, HOST_APP_PRO, PORT_APP_PRO,
                    CREDENTIALS_PRO, DEBUG_PRO, METHODS, ALLOWED_HOSTS_PRO,DSN,HASH_KEY
                )

                app.debug = DEBUG_PRO
                host = HOST_APP_PRO
                port = PORT_APP_PRO
                dsn = DSN 
                SECRET_KEY = HASH_KEY

                # Production Middleware

                app.add_middleware(CORSMiddleware, allow_headers=HEADERS, allow_origins=ORIGINS_PRO,
                                   allow_credentials=CREDENTIALS_PRO, allow_methods=METHODS)
                app.add_middleware(TrustedHostMiddleware,
                                   allowed_hosts=ALLOWED_HOSTS_PRO)
                app.add_middleware(GZipMiddleware, minimum_size=1000)

            else:
                from settings.dev import (
                    CORSMiddleware, ORIGINS, CREDENTIALS, METHODS, HEADERS,
                    HOST_APP, PORT_APP, DEBUG,DSN_DEV,HASH_KEY_DEV)

                app.debug = DEBUG
                host = HOST_APP
                port = PORT_APP
                dsn = DSN_DEV
                SECRET_KEY = HASH_KEY_DEV

                # Develop middleware

                app.add_middleware(CORSMiddleware, allow_headers=HEADERS, allow_origins=ORIGINS,
                                   allow_credentials=CREDENTIALS, allow_methods=METHODS)

            if selection.host:
                host = selection.host
            if selection.port:
                port = selection.port
            
            print(f"[+] Runing server debug mode {app.debug}")
            setDatabaseConn(dsn)
            server.run('main:app', host=host,
                       port=port, reload=True, workers=4)

        case ('tests'):
            # Test mode

            if selection.file:
                print("[+] Runing tests...")
                subprocess.run(['python -m unittest', f"{selection.file}"])
                
            else:
                print("[-] File not found")

        case ('database'):
            # Database mode

            print(
                f"[-] This backend don't manage database only fiber-backend {selection.argument}")

        case ('shell'):
            # Shell mode

            print("[+] Runing shell...")
            if selection.file:
                subprocess.run(['python', f"< {selection.file}"])
                
            subprocess.run(["python"])

        case _:
            # Default

            print("[-] Command not reconized")
            exit(0)