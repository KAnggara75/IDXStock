import {
	PrismaClientInitializationError,
	PrismaClientKnownRequestError,
	PrismaClientRustPanicError,
	PrismaClientUnknownRequestError,
	PrismaClientValidationError,
} from "@prisma/client/runtime/library";
import { ZodError } from "zod";
import { log } from "../config/logger";
import { JsonUtils } from "../utils/jsonUtils";
import { HTTPException } from "hono/http-exception";

import type { Context } from "hono";

export async function errorHandler(err: unknown, c: Context) {
	if (err instanceof HTTPException) {
		const response: unknown = JsonUtils.safeParseJSON(err.message);
		return c.json({ errors: response }, err.status);
	}

	if (err instanceof ZodError) {
		return c.json({ errors: JsonUtils.safeParseJSON(err.message) }, 400);
	}

	if (
		err instanceof PrismaClientKnownRequestError ||
		err instanceof PrismaClientUnknownRequestError ||
		err instanceof PrismaClientRustPanicError ||
		err instanceof PrismaClientValidationError
	) {
		log.error(err.constructor.name);
		return c.json({ errors: JsonUtils.safeParseJSON(err.message) }, 503);
	}

	if (err instanceof PrismaClientInitializationError) {
		log.error("PrismaClientInitializationError");
		return c.json(
			{
				errors: "Can't reach database server",
				code: err.errorCode,
			},
			503
		);
	}

	log.error(`Unknown error: ${(err as Error).message}`);
	return c.json({ errors: (err as Error).message }, 500);
}
