import app from "./index.ts";
import { log } from "./config/logger";
import type { Logger } from "winston";
import { prismaClient } from "./config/database";
import { serve, type ServerType } from "@hono/node-server";

const port: number = Number(Bun.env.API_PORT ?? 3000);

const server: ServerType = serve({
	port,
	fetch: app.fetch,
});

const events: string[] = ["uncaughtException", "SIGINT", "SIGTERM"];

events.forEach((eventName: string): void => {
	process.on(eventName, (...args): void => {
		gracefulShutdown().then(
			(r: number): Logger => log.warn(`gracefulShutdown with ${r}`)
		);
		log.warn(
			`${eventName} was called with args: ${args.join(",")}, closing HTTP server`
		);
	});
});

async function gracefulShutdown(): Promise<number> {
	let status: number = 0;
	log.warn("Shutting down gracefully...");
	await prismaClient.$disconnect();

	server.close((): never => {
		log.warn("HTTP server closed");
		process.exit(0);
	});

	setTimeout((): never => {
		log.error("Could not close connections in time, forcefully shutting down");
		status = 0;
		process.exit(1);
	}, 5000);

	return status;
}
