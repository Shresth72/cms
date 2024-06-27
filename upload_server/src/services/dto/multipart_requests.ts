import { IsNumber, IsString } from "class-validator";

class UploadRequest {
  @IsString()
  filename: string;
}

class UploadChunkRequest {
  @IsString()
  filename: string;
  @IsNumber()
  chunkIndex: number;
  uploadId: string;
}

class CompleteUploadRequest {
  @IsString()
  filename: string;
  totalChunks: number;
  uploadId: string;
  title: string;
  description: string;
  author: string;
}

class UploadToDb {
  //
}

export { UploadRequest, UploadChunkRequest, CompleteUploadRequest, UploadToDb };
