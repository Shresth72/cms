import mongoose from "mongoose";

const Schema = mongoose.Schema;

const MetadataSchema = new Schema({
    title: {
        type: String,
        required: true,
    },
    desc: {
        type: String,
        required: true,
    },
    url: {
        type: String,
        required: true,
    },
    author: {
        type: String,
        required: false,
    },
    createdAt: {
        type: Date,
        default: Date.now,
    },
    updatedAt: {
        type: Date,
        default: Date.now,
    },
    // todo: Add more
});

export default mongoose.model("metadata", MetadataSchema);
