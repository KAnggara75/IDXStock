import { ZodError, ZodSchema } from "zod";
import { validator } from "hono/validator";

export function validateWithSchema<T>(type: "json", schema: ZodSchema<T>) {
	return validator(type, (value) => {
		const result = schema.safeParse(value);
		if (!result.success) {
			throw new ZodError(result.error.issues);
		}
		return result.data;
	});
}
