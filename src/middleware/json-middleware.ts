import type { MiddlewareHandler } from "hono";
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
