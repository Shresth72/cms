import AWS from "aws-sdk";

class S3 {
  s3: AWS.S3;
  bucketName: string;
}

export function initializeS3(): S3 {
  const s3 = new AWS.S3({
    accessKeyId: process.env.AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
    region: "ap-south-1",
  });

  const bucketName = process.env.BUCKET_NAME;

  return { s3, bucketName };
}
