import { describe, expect, test } from "bun:test";
import { FileValidation } from "../src/validation/file-validation";
import { ZodError } from "zod";

describe("FileValidation.DOCUMENT_CHECK", (): void => {
	const path = "./test/Stock Summary-12345678.xlsx";
	const blob = Bun.file(path);

	it("Should return true if the file name does not exist", async (): Promise<void> => {
		const file = new File([blob], "Stock Summary-12345678.xlsx", {
			type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		});
		expect(() => FileValidation.DOCUMENT_CHECK.parse(file)).not.toThrow();
	});

	it("Should return true if the file name does not exist", async (): Promise<void> => {
		const file = new File([blob], "Stock-12345678.xlsx", {
			type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		});
		expect(() => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			"Invalid document file name"
		);
	});

	it("Invalid file type should fail", async (): Promise<void> => {
		const file = new File([blob], "Stock Summary-12345678.xlsx", {
			type: "text/plain",
		});
		expect(() => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			"Invalid document file type"
		);
	});
});
