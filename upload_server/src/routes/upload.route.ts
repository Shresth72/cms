import express from "express";
import MultipartUploadService from "../services/multipart_upload.service";
import multer from "multer";

const upload = multer();

const router = express.Router();

const multipartUploadService = new MultipartUploadService();
// TODO: Add routes for smaller file uploads

// start upload route
router.post("/init", upload.none(), multipartUploadService.initUpload);

// upload each chunk
router.post("/", upload.single("chunk"), multipartUploadService.uploadChunk);

// complete multipart upload
router.post("/complete", multipartUploadService.completeUpload);

// upload metadata to mongo
// router.post("/uploadDb", multipartUploadService.uploadToDb);

export default router;
