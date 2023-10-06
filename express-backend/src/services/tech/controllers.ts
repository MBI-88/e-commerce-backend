import { Request, Response } from "express";


export async function GetTech(req:Request,res:Response) {
    res.send("Hola Tech")
}