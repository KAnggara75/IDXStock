import { Hono } from "hono";
import { ZodError } from "zod";
import pkg from "../package.json";
import { log } from "./config/logger";
import { prismaClient } from "./config/database";
import { HTTPException } from "hono/http-exception";
import { serve, type ServerType } from "@hono/node-server";
import { userController } from "./controller/user-controller";
import { contactController } from "./controller/contact-controller";
import { addressController } from "./controller/address-controller";

const port: number = Number(Bun.env.API_PORT ?? 3030);

export const app = new Hono().basePath("/api");

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

app.route("/", userController);
app.route("/", contactController);
app.route("/", addressController);

app.notFound((c) => {
	log.info(`Not Found: ${c.req.url}`);
	return c.json({ errors: "Not Found" }, 404);
});

app.onError(async (err, c) => {
	if (err instanceof HTTPException) {
		return c.json({ errors: err.message }, err.status);
	} else if (err instanceof ZodError) {
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
		gracefulShutdown();
		log.info(`${eventName} was called with args : ${args.join(",")}`);
		log.info(`${eventName} signal received: closing HTTP server`);
	});
});

async function gracefulShutdown(): Promise<void> {
	log.info("Shutting down gracefully...");
	await prismaClient.$disconnect();

	server.close(() => {
		log.info("HTTP server closed");
		// Close any other connections or resources here
		process.exit(0);
	});

	// Force close the server after 5 seconds
	setTimeout(() => {
		console.error(
			"Could not close connections in time, forcefully shutting down"
		);
		process.exit(1);
	}, 5000);
}

export default server;
