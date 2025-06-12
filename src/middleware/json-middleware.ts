import type { MiddlewareHandler } from "hono";
import { log } from "../config/logger.ts";
import { type CustomError, toErrorDetail } from "../model/errors-model.ts";
import { HTTPException } from "hono/http-exception";

const methodsWithBody = ["POST", "PUT", "PATCH"];

function isMethodsWithBody(method: string): boolean {
	return methodsWithBody.includes(method);
}

export const jsonMiddleware: MiddlewareHandler = async (c, next) => {
	console.warn("jsonMiddleware");

	if (!isMethodsWithBody(c.req.method.toUpperCase())) {
		return next();
	}

	if (c.req.method === "GET" || c.req.method === "DELETE") {
		return next();
	}

	if (c.req.path === "/api/convert") {
		return next();
	}

	if (c.req.path.includes("/api/stocks/idx")) {
		return next();
	}

	const contentType = c.req.header("content-type");

	if (!contentType) {
		return c.json(
			{
				errors: "Invalid Header, Content-Type Not Found",
			},
			400
		);
	}

	if (!contentType.includes("application/json")) {
		return c.json(
			{
				errors: `Invalid Content-Type, expected JSON, received ${contentType}`,
			},
			400
		);
	}

	await next();
};

export const validateJsonBody: MiddlewareHandler = async (c, next) => {
	if (!isMethodsWithBody(c.req.method.toUpperCase())) {
		return next();
	}

	const contentType = c.req.header("content-type");
	if (contentType?.includes("application/json")) {
		log.debug("validateJsonBody");
		try {
			const body = await c.req.json();
			c.set("parsedBody", body);
			await next();
		} catch {
			const errorPayload: CustomError[] = await toErrorDetail(
				"bad_request",
				"Invalid or missing JSON body"
			);
			throw new HTTPException(400, {
				message: JSON.stringify(errorPayload),
			});
		}
	}
};
