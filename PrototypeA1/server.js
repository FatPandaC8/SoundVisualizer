const express = require("express");

const app = express();
const RANGE = [40, 80, 120, 180, 300];

// allow JSON FFT frames
app.use(express.json({ limit: "2mb" }));

// serve frontend
app.use(express.static("public"));

app.post("/fft", (req, res) => {
  const { freq } = req.body;

  console.log(freq);

  res.sendStatus(200);
});

app.listen(3000, () => {
  console.log("Server running at http://localhost:3000");
});