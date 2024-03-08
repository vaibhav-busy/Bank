import express from 'express';
import nodeBankRoutes from "./routes_node/bankRoute.js";
import bodyParser from "body-parser";

const app=express();
app.use(bodyParser.json());

const PORT=3000;

// app.get('/',(req,res)=>{
//     res.status.json("Successfully Connected to node")
// })

app.get('/', (req, res) => {
    res.status(200).json({ message: "Successfully Connected to node" });
  });
  

app.use('/bank',nodeBankRoutes)
// app.use('/branch',nodeBranchRoutes)


app.listen(PORT, () => {
    console.log(`Node server is listening on PORT ${PORT}`);
  });
