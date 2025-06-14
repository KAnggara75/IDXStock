import { log } from "./logger";
import { Prisma, PrismaClient } from "@prisma/client";
import { isSensitiveQuery } from "../helpers/sensitive-helper";

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
	const sensitive = isSensitiveQuery(e.query);
	log.debug(`SQL|Query: ${sensitive ? "[SENSITIVE QUERY HIDDEN]" : e.query}`);
	log.debug(`SQL|Params: ${sensitive ? "[REDACTED]" : e.params}`);
	log.debug(`SQL|Duration: ${e.duration} ms`);
});
