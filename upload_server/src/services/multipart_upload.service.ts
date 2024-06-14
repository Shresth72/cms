import express from "express";

import MetadataRepository from "../database/repository/metadata.repository";

class MultipartUploadService {
  constructor() {
    this.repository = new MetadataRepository();
    // this.initUpload = this.initUpload.bind(this);
  }

  async initUpload(req, res) {}

  async uploadChunk(req, res) {}

  async completeUpload(req, res) {}

  async uploadToDb(req, res) {}
}

export default MultipartUploadService;
