import { log } from "./logger";
import { Prisma, PrismaClient } from "@prisma/client";
import { isSensitiveQuery } from "../helpers/sensitive-helper";

const logLevels: Prisma.LogLevel[] = ["query", "error", "info", "warn"];

export const prismaClient = new PrismaClient({
	log: logLevels.map(
		(level: Prisma.LogLevel): { emit: "event"; level: Prisma.LogLevel } => ({
			emit: "event",
			level,
		})
	),
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
	const sensitive: boolean = isSensitiveQuery(e.query);
	log.debug(`SQL|Query: ${sensitive ? "[SENSITIVE QUERY HIDDEN]" : e.query}`);
	log.debug(`SQL|Params: ${sensitive ? "[REDACTED]" : e.params}`);
	log.debug(`SQL|Duration: ${e.duration} ms`);
});
