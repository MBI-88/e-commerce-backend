import { Router } from "express";
import { GetHome } from "./controllers";


const news = Router()

news.use("/home",GetHome)


export default news