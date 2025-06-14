// file: middleware/requestLogger.ts
import { v4 as uuid } from "uuid";
import { log } from "../config/logger";
import { requestContext } from "../helpers/request-helper";
import type { Context, MiddlewareHandler, Next } from "hono";

export const requestLogger: MiddlewareHandler = async (
	c: Context,
	next: Next
): Promise<void> => {
	const id: string = uuid();
	const url: string = c.req.url;
	const method: string = c.req.method;
	const start: number = Date.now();

	const ip: string =
		c.req.header("x-forwarded-for") ||
		c.req.header("cf-connecting-ip") ||
		c.req.header("x-real-ip") ||
		"unknown";

	return requestContext.run({ requestId: id }, async () => {
		c.set("requestId", id);
		c.header("x-request-id", id);
		log.info(`➡️ ${method} ${url} from ${ip}`);

		try {
			await next();
		} finally {
			const duration: number = Date.now() - start;
			const status: number = c.res.status;
			log.info(`✅ ${method} ${url} ${status} - ${duration}ms`);
		}
	});
};
