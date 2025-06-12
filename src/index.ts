import { Hono } from "hono";
import { ZodError } from "zod";
import pkg from "../package.json";
import { log } from "./config/logger";
import {
	PrismaClientInitializationError,
	PrismaClientKnownRequestError,
	PrismaClientRustPanicError,
	PrismaClientUnknownRequestError,
	PrismaClientValidationError,
} from "@prisma/client/runtime/library";
import { prismaClient } from "./config/database";
import { HTTPException } from "hono/http-exception";
import { serve, type ServerType } from "@hono/node-server";
import { userController } from "./controller/user-controller";
import { authController } from "./controller/auth-controller";
import { jsonMiddleware, validateJsonBody } from "./middleware/json-middleware";
import { stockController } from "./controller/stock-controller";
import { converterController } from "./controller/convert-controller.ts";
import { authMiddleware } from "./middleware/auth-middleware.ts";
import { JsonUtils } from "./utils/jsonUtils.ts";

const port: number = Number(Bun.env.API_PORT ?? 3000);

export const app = new Hono().basePath("/api");

app.use(validateJsonBody);

app.get("/version", (c) => {
	log.debug("get /version");
	return c.json(
		{
			message: pkg.description ?? "KAnggara Web APP",
			version: Bun.env.APP_VERSION ?? pkg.version ?? "0.0.1",
			stability: Bun.env.APP_STAB ?? "Developer-Preview",
		},
		200
	);
});
app.route("/auth", authController);

app.use(jsonMiddleware);
app.use(authMiddleware);

app.route("/users", userController);
app.route("/stocks", stockController);
app.route("/convert", converterController);

app.notFound((c) => {
	log.error(`Not Found: ${c.req.path}`);
	return c.json({ errors: "Not Found" }, 404);
});

app.onError(async (err, c) => {
	if (err instanceof HTTPException) {
		const response = JsonUtils.safeParseJSON(err.message);
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

	log.error("Unknown error: " + (err as Error).message);
	return c.json({ errors: (err as Error).message }, 500);
});

const server: ServerType = serve({
	port: port,
	fetch: app.fetch,
});

const events = ["uncaughtException", "SIGINT", "SIGTERM"];

events.forEach((eventName) => {
	process.on(eventName, (...args) => {
		gracefulShutdown().then((r) => log.warn(`gracefulShutdown with ${r}`));
		log.warn(
			`${eventName} was called with args: ${args.join(",")} closing HTTP server`
		);
	});
});

async function gracefulShutdown(): Promise<number> {
	let status: number = 0;
	log.warn("Shutting down gracefully...");
	await prismaClient.$disconnect();

	// Close any other connections or resources here
	server.close(() => {
		log.warn("HTTP server closed");
		process.exit(0);
	});

	// Force close the server after 5 seconds
	setTimeout(() => {
		log.error("Could not close connections in time, forcefully shutting down");
		status = 0;
		process.exit(1);
	}, 5000);

	return status;
}

export default server;
