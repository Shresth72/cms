import { Request, Response } from "express";
import { validate } from "class-validator";
import { plainToInstance } from "class-transformer";
import {
  UploadRequest,
  UploadChunkRequest,
  CompleteUploadRequest,
  UploadToDb,
} from "./dto/multipart_requests";
import AWS from "aws-sdk";
import dotenv from "dotenv";

dotenv.config();

import MetadataRepository from "../database/repository/metadata.repository";

class MultipartUploadService {
  repository: MetadataRepository;

  constructor() {
    this.repository = new MetadataRepository();
    // this.initUpload = this.initUpload.bind(this);
  }

  async initUpload(req: Request, res: Response) {
    try {
      const uploadRequest = plainToInstance(UploadRequest, req.body);
      const errors = await validate(uploadRequest);

      if (errors.length > 0) {
        return res
          .status(400)
          .json({ errors: errors.map((err) => err.toString()) });
      }

      const { filename } = uploadRequest;

      // AWS S3 INIT
      // TODO: Setup AWS ENV from terraform
      const s3 = new AWS.S3({
        accessKeyId: process.env.AWS_ACCESS_KEY_ID,
        secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
        region: "ap-south-1",
      });

      // TODO: To be passed from terraform
      const bucketName = process.env.BUCKET_NAME;

      const s3Params = {
        Bucket: bucketName,
        Key: filename,
        ContentType: "video/mp4",
      };

      const multipartParams = await s3
        .createMultipartUpload(s3Params)
        .promise();
      const uploadId = multipartParams.UploadId;

      res.status(200).json({ uploadId });
    } catch (err) {
      console.log("Error initializing upload", err);
      res.status(500).send("Upload initialization failed");
    }
  }

  async uploadChunk(req: Request, res: Response) { }

  async completeUpload(req: Request, res: Response) { }

  // abort

  async uploadToDb(req: Request, res: Response) { }
}

export default MultipartUploadService;
