export interface CustomError {
	validation?: string;
	code: string;
	expected?: string;
	received?: string;
	path?: (string | number)[];
	message?: string;
}

export async function toErrorDetail(
	code: string,
	message: string,
	path?: (string | number)[],
	validation?: string
): Promise<CustomError[]> {
	return [{ validation: validation, code: code, path: path, message: message }];
}
