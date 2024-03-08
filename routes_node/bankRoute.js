import axios from "axios";
import express from "express";


const router = express.Router()

const headers = {
    'Content-Type': 'application/json'
};


router.get("/" , async(req ,res) => {
    await axios.get("http://localhost:8080/bank", {headers}).then(
        response => {
            const banks = response.data
            console.log(banks)
            res.status(200).json({"message" : "Bank fetched successfully"})
        }
    ).catch(err => {
        res.status(400).json({"error" : err.response.data.error})
    })
})


router.post("/" , async(req ,res) => {
    await axios.post("http://localhost:8080/bank", JSON.stringify(req.body), {headers}).then(
        response => {
            console.log(response.data)
            res.status(200).json({"message" : "Bank created successfully"})
        }
    ).catch(err => {
        res.status(400).json({"error" : err.response.data.error})
    })
})


router.patch("/" , async(req ,res) => {
    await axios.patch("http://localhost:8080/bank",JSON.stringify(req.body), {headers}).then(
        response => {
            console.log(response.data)
            res.status(200).json({"message" : "Bank updated successfully"})
        }
    ).catch(err => {
        res.status(400).json({"error" : err.response.data.error})
    })
})

router.delete("/:id" , async(req ,res) => {
    bankId = req.params.id
    await axios.delete(`http://localhost:8080/bank/${bankId}`, {headers}).then(
        response => {
            console.log(response.data)
            res.status(200).json({"message" : "Bank deleted successfully"})
        }
    ).catch(err => {
        res.status(400).json({"error" : err.response.data.error})
    })
})

router.get("/:id" , async(req ,res) => {
    bankId = req.params.id
    await axios.get(`http://localhost:8080/bank/${bankId}`, {headers}).then(
        response => {
            console.log(response.data)
            res.status(200).json({"message" : "Bank details fetched successfully"})
        }
    ).catch(err => {
        res.status(400).json({"error" : err.response.data.error})
    })
})



export default router;
