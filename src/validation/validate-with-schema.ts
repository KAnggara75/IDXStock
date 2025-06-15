import { type SafeParseReturnType, ZodError, ZodSchema } from "zod";
import { validator } from "hono/validator";

export function validateWithSchema<T>(
	schema: ZodSchema<T>,
	type: "json" = "json"
) {
	return validator(type, (value: unknown): T => {
		const result: SafeParseReturnType<T, T> = schema.safeParse(value);
		if (!result.success) {
			throw new ZodError(result.error.issues);
		}
		return result.data;
	});
}
