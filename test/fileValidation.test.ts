import { ZodError } from "zod";
import type { BunFile } from "bun";
import { describe, expect, it } from "bun:test";
import { FileValidation } from "../src/validation/file-validation";

describe("FileValidation.DOCUMENT_CHECK", (): void => {
	const path = "./test/Stock Summary-12345678.xlsx";
	const blob: BunFile = Bun.file(path);

	it("Should return true if the file name does not exist", async (): Promise<void> => {
		const file = new File([blob], "Stock Summary-12345678.xlsx", {
			type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		});
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).not.toThrow();
	});

	it("Should return true if the file name does not exist", async (): Promise<void> => {
		const file = new File([blob], "Stock-12345678.xlsx", {
			type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		});
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			"Invalid document file name"
		);
	});

	it("Invalid file type should fail", async (): Promise<void> => {
		const file = new File([blob], "Stock Summary-12345678.xlsx", {
			type: "text/plain",
		});
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			"Invalid document file type"
		);
	});
});

describe("FileValidation.DOCUMENT_CHECK", (): void => {
	it("should pass with valid XLSX file", (): void => {
		const file: File = mockFile(
			"Stock Summary-June.xlsx",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		);
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).not.toThrow();
	});

	it("should fail with invalid file type", (): void => {
		const file: File = mockFile("Stock Summary-June.xlsx", "application/pdf");
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			ZodError
		);
	});

	it("should fail with invalid file name", (): void => {
		const file: File = mockFile(
			"Report-June.xlsx",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		);
		expect((): void => FileValidation.DOCUMENT_CHECK.parse(file)).toThrow(
			ZodError
		);
	});

	it("should fail if not a File instance", (): void => {
		expect((): void => FileValidation.DOCUMENT_CHECK.parse({})).toThrow(
			ZodError
		);
	});
});

function mockFile(name: string, type: string): File {
	return new File(["dummy content"], name, { type });
}
