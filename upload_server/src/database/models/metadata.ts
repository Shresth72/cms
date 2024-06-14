import mongoose from "mongoose";

const Schema = mongoose.Schema;

const MetadataSchema = new Schema({
  title: String,
  desc: String,
  url: String,
  // todo: Add more
});

export default mongoose.model("metadata", MetadataSchema);
