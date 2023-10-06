import express from "express";
import router from "./services/router/router";
import * as readline from 'readline'



const app = express()

const port = 3001


app.use(router)





app.listen(port, () =>{
    console.log("D'Todo express backend");
    
})