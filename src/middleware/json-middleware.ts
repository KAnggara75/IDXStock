import type { MiddlewareHandler } from "hono";

export const jsonMiddleware: MiddlewareHandler = async (c, next) => {
	if (c.req.method === "GET") {
		return next();
	}

	const contentType = c.req.header("content-type");

	if (!contentType || !contentType.includes("application/json")) {
		return c.json(
			{
				errors: `Invalid Content-Type, expected JSON, received ${contentType}`,
			},
			400
		);
	}

	await next();
};