import { ErrorHandlingMiddlewareFunction, ObjectId } from "mongoose";
import MetadataModel from "../models/metadata";
// TODO: Add error handler

class MetadataRepository {
    async CreateFile(data: MetadataRepository) {
        try {
            const metadata = new MetadataModel(data);
            await metadata.save();
            return metadata;
        } catch (err) {
            console.log("Error creating metadata file");
            throw new Error("Could not create metadata file.");
        }
    }

    async FindFile(id: ObjectId) {
        try {
            const metadata = await MetadataModel.findById(id);
        } catch (err) {
            console.error("Error finding metadata file.", err);
            throw new Error("Could not find metadata file.");
        }
    }

    async FindFileByName(search: string) {
        try {
            const metadata = await MetadataModel.find({ title: search });
        } catch (err) {
            console.error("Error finding metadata file.", err);
            throw new Error("Could not find metadata file by name.");
        }
    }

    async UpdateFile(id: ObjectId, data: MetadataRepository) {
        try {
            const metadata = await MetadataModel.findByIdAndUpdate(id, data);
        } catch (err) {
            console.error("Error Updating metadata files.", err);
            throw new Error("Could not update metadata files");
        }
    }

    async DeleteFile(id: ObjectId) {
        try {
            const metadata = await MetadataModel.findByIdAndDelete(id);
            // console.log(metadata);
        } catch (err) {
            console.error("Error Deleting metadata files.", err);
            throw new Error("Could not delete metadata files.");
        }
    }
}

export default MetadataRepository;
