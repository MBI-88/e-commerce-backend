import { Response,Request } from "express";


export async function GetHome(req:Request,res:Response) {
    res.send("Hola Home")
}