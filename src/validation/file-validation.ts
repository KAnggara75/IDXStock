import { z, ZodType } from "zod";

export class FileValidation {
	static readonly DOCUMENT_CHECK: ZodType = z
		.instanceof(File)
		.refine(
			(file) =>
				[
					"application/vnd.ms-excel",
					"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				].includes(file.type),
			{ message: "Invalid document file type" }
		)
		.refine((file) => file.name.startsWith("Stock Summary-"), {
			message: "Invalid document file name",
		});
}
