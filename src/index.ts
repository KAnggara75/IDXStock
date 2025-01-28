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
import { jsonMiddleware } from "./middleware/json-middleware";
import { stockController } from "./controller/stock-controller";
import { converterController } from "./controller/convert-controller.ts";

const port: number = Number(Bun.env.API_PORT ?? 3030);

export const app = new Hono().basePath("/api");

authController.use(jsonMiddleware);

app.get("/", (c) => {
	return c.json(
		{
			message: pkg.description ?? "KAnggara Web APP",
			version: Bun.env.APP_VERSION ?? pkg.version ?? "0.0.1",
			stability: Bun.env.APP_STAB ?? "Developer-Preview",
		},
		200
	);
});

app.route("/", authController);
app.route("/", userController);
app.route("/", converterController);
app.route("/", stockController);

app.notFound((c) => {
	log.error(`Not Found: ${c.req.path}`);
	return c.json({ errors: "Not Found" }, 404);
});

app.onError(async (err, c) => {
	if (err instanceof HTTPException) {
		return c.json({ errors: err.message }, err.status);
	} else if (err instanceof ZodError) {
		return c.json({ errors: JSON.parse(err.message) }, 400);
	} else if (err instanceof PrismaClientKnownRequestError) {
		log.error("PrismaClientKnownRequestError");
		return c.json({ errors: err.message }, 400);
	} else if (err instanceof PrismaClientUnknownRequestError) {
		log.error("PrismaClientUnknownRequestError");
		return c.json({ errors: JSON.parse(err.message) }, 400);
	} else if (err instanceof PrismaClientRustPanicError) {
		log.error("PrismaClientRustPanicError");
		return c.json({ errors: JSON.parse(err.message) }, 400);
	} else if (err instanceof PrismaClientInitializationError) {
		log.error("PrismaClientInitializationError");
		return c.json(
			{
				errors: "Can't reach database server",
				code: err.errorCode,
			},
			503
		);
	} else if (err instanceof PrismaClientValidationError) {
		log.error("PrismaClientValidationError");
		return c.json({ errors: JSON.parse(err.message) }, 400);
	} else {
		return c.json({ errors: err.message }, 500);
	}
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
		console.error(
			"Could not close connections in time, forcefully shutting down"
		);
		status = 0;
		process.exit(1);
	}, 5000);

	return status;
}

export default server;
