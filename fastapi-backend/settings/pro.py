from .dev import (get_config, HEADERS, GZipMiddleware,
                    TrustedHostMiddleware, CORSMiddleware)


# Settings

settings = get_config()

DEBUG_PRO = settings.DEBUG

# Origins

ORIGINS_PRO = [origin for origin in settings.ORIGINS.split(',')]

# HOST APP

HOST_APP_PRO = settings.HOST_APP

# PORT APP

PORT_APP_PRO = settings.PORT_APP

# Credentials

CREDENTIALS_PRO = settings.CREDENTIALS

# Allowed methods

METHODS = [method for method in settings.METHODS.split(',')]

# Allowed_hosts 

ALLOWED_HOSTS_PRO = [host for host in settings.ALLOW_HOSTS.split(',')]

# Database

DSN = settings.DSN

# Hash key 

HASH_KEY = settings.HASHKEY
