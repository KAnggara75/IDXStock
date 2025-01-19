import type { MiddlewareHandler } from "hono";

export const jsonMiddleware: MiddlewareHandler = async (c, next) => {
	if (c.req.method === "GET" || c.req.method === "DELETE") {
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
