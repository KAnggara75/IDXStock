import { log } from "./logger";
import { Prisma, PrismaClient } from "@prisma/client";

export const prismaClient = new PrismaClient({
	log: [
		{
			emit: "event",
			level: "query",
		},
		{
			emit: "event",
			level: "error",
		},
		{
			emit: "event",
			level: "info",
		},
		{
			emit: "event",
			level: "warn",
		},
	],
});
prismaClient.$on("error", (e: Prisma.LogEvent): void => {
	log.error(e);
});
prismaClient.$on("warn", (e: Prisma.LogEvent): void => {
	log.warn(e);
});
prismaClient.$on("info", (e: Prisma.LogEvent): void => {
	log.info(e);
});
prismaClient.$on("query", (e: Prisma.QueryEvent): void => {
	log.debug(`SQL|Query: ${e.query}`);
	log.debug(`SQL|Params: ${e.params}`);
	log.debug(`SQL|Duration: ${e.duration} ms`);
});
