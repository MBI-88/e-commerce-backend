from datetime import datetime
from fastapi import HTTPException, status, Depends
from jose import jwt, JWTError
from fastapi.security import OAuth2PasswordBearer


oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/users/signin")
SECRET_KEY:str = ""


def verify_access_token(token: str) -> dict:
    try:
        data = jwt.decode(token, SECRET_KEY, algorithms="HS256")
        expire = data.get("exp")

        if expire is None:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST, detail="No access token supplied"
            )

        if datetime.utcnow() > datetime.utcfromtimestamp(expire):
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN, detail="Token expired!"
            )

        return data
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST, detail="Invalid token")
    

async def authenticate(token: str = Depends(oauth2_scheme)) -> dict:
    if not token:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN, detail="Sign in for access")

    user = verify_access_token(token)
    return user