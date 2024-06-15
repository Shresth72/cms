import { Request, Response } from "express";
import { validate } from "class-validator";
import { plainToInstance } from "class-transformer";
import {
  UploadRequest,
  UploadChunkRequest,
  CompleteUploadRequest,
  UploadToDb,
} from "./dto/multipart_requests";
import dotenv from "dotenv";
import { initializeS3 } from "../utils";

// const upload = multer();

dotenv.config();

import MetadataRepository from "../database/repository/metadata.repository";

class MultipartUploadService {
  repository: MetadataRepository;

  constructor() {
    this.repository = new MetadataRepository();
    // this.initUpload = this.initUpload.bind(this);
  }

  async initUpload(req: Request, res: Response): Promise<Response> {
    try {
      const uploadRequest = plainToInstance(UploadRequest, req.body);
      const errors = await validate(uploadRequest);

      if (errors.length > 0) {
        return res
          .status(400)
          .json({ errors: errors.map((err) => err.toString()) });
      }

      const { filename } = uploadRequest;
      const { s3, bucketName } = initializeS3();

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
      console.log("Error initializing upload: ", err);
      res.status(500).send("Upload initialization failed");
    }
  }

  async uploadChunk(req: Request, res: Response): Promise<Response> {
    try {
      const uploadChunkRequest = plainToInstance(UploadChunkRequest, req.body);
      const errors = await validate(uploadChunkRequest);

      if (errors.length > 0) {
        return res
          .status(400)
          .json({ errors: errors.map((err) => err.toString()) });
      }

      const { filename, chunkIndex, uploadId } = uploadChunkRequest;
      const { s3, bucketName } = initializeS3();

      // Added through middleware
      if (!req.file) {
        return res.status(400).json({ error: "No file uploaded" });
      }

      const partParams = {
        Bucket: bucketName,
        Key: filename,
        UploadId: uploadId,
        PartNumber: chunkIndex + 1,
        Body: req.file.buffer,
      };

      await s3.uploadPart(partParams).promise();
      res.status(200).json({ success: true });
    } catch (err) {
      console.log("Error uploading chunks: ", err);
      res.status(500).send("Could not upload chunk");
    }
  }

  async completeUpload(req: Request, res: Response): Promise<Response> {
    const completeUploadRequest = plainToInstance(
      CompleteUploadRequest,
      req.body,
    );
    const errors = await validate(completeUploadRequest);

    if (errors.length > 0) {
      return res
        .status(400)
        .json({ errors: errors.map((err) => err.toString()) });
    }

    const { filename, totalChunks, uploadId, title, description, author } =
      completeUploadRequest;

    const uploadedParts = [];

    for (let i = 0; i < totalChunks; i++) {
      uploadedParts.push({
        PartNumber: i + 1,
        ETag: req.body[`parts${i + 1}`],
      });
    }

    const { s3, bucketName } = initializeS3();

    const completeParams = {
      Bucket: bucketName,
      Key: filename,
      UploadId: uploadId,
      MultiPartUpload: {
        Parts: uploadedParts,
      },
    };

    const uploadResult = await s3
      .completeMultipartUpload(completeParams)
      .promise();

    // await addVideoDetailsToDB(title, description, author, uploadResult.Location);
    // Call Kafka for encoding
    return res.status(200).json({ message: "Uploaded Successfully!" });
  }

  // abort

  // async uploadToDb(req: Request, res: Response): Promise<Response> { }
}

export default MultipartUploadService;
