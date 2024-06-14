import express from "express";
import cors from "cors";
import dotenv from "dotenv";

import databaseConnection from "./database/connections";
import uploadRouter from "./routes/upload.route";
import kafkaRouter from "./routes/kafka.route";

const port = process.env.PORT || 8080;

<<<<<<< Updated upstream
const StartServer = async () => {
  dotenv.config();

  const app = express();

  await databaseConnection();

  // TODO: Add Cors for Next
  app.use(
    cors({
      allowedHeaders: ["*"],
      origin: "*",
    }),
  );

  app.use(express.json());
  app.use(
    express.urlencoded({
      extended: true,
    }),
  );

  app.use("/upload", uploadRouter);
  app.use("/publish", kafkaRouter);

  // TODO: Add Error Handling
  app.listen(port, () => {
    console.log(`Server listening on a ${port}`);
  });
};

StartServer();
=======
const app = express();
app.use(
    cors({
        allowedHeaders: ["*"],
        origin: "*",
    })
);

app.get("/", (req, res) => {
    res.send("Hello ");
});

app.listen(port, () => {
    console.log(`Server listening on ${port}`);
});
>>>>>>> Stashed changes
