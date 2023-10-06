import { Router } from "express";
import { GetTech } from "./controllers";

const tech = Router()


tech.use("/home",GetTech)

export default tech