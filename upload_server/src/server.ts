import express from "express";
import cors from "cors";
import dotenv from "dotenv";

dotenv.config();

const port = process.env.PORT || 8080;

const app = express();
app.use(
  cors({
    allowedHeaders: ["*"],
    origin: "*",
  }),
);

app.get("/", (req, res) => {
  res.send("Hello ");
});

app.listen(port, () => {
  console.log(`Server listening on ${port}`);
});
