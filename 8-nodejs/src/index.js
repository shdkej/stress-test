"use strict";

const express = require("express");

// Constants
const PORT = process.env.PORT || 8080;
const HOST = "0.0.0.0";

// App
const app = express();

// API
app.get("/", (req, res) => {
    res.set("Content-Type", "application/json");
    res.send("1");
});

app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);
