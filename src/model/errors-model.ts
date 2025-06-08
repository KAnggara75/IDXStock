export interface CustomError {
	code: string;
	expected?: string;
	received?: string;
	path: string[];
	message: string;
}

export async function toErrorDetail(
	code: string,
	path: string[],
	message: string
): Promise<CustomError[]> {
	return [
		{
			code: code,
			path: path,
			message: message,
		},
	];
}
