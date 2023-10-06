import { Router } from "express";
import news from "../news/routes";
import tech from "../tech/routes";



const router = Router()


router.use("/news",news)
router.use("/tech",tech)


export default router