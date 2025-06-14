import { Format, type TransformableInfo } from "logform";
import { getRequestId } from "../helpers/request-helper";
import { createLogger, format, type Logger, transports } from "winston";

const { combine, timestamp, printf, colorize, errors } = format;

const logsFormat: Format = printf(
	({ level, message, timestamp, stack }: TransformableInfo): string => {
		const logMessage: unknown = stack || message;
		const requestId = getRequestId() ?? " ".padEnd(36, " ");
		return `${timestamp}|[${level.toUpperCase().padEnd(5)}]|[${requestId}]|${logMessage}`;
	}
);

const logLevel: string = Bun.env.LOG_LEVEL ?? "warn";
const todayDate: string = new Date().toISOString().slice(0, 10);
const logFolder: string = `./logs/${todayDate.replace(/-/g, "")}/${todayDate}`;

const myTransports: (
	| transports.ConsoleTransportInstance
	| transports.FileTransportInstance
)[] = [
	new transports.File({
		filename: `${logFolder}-app.log`,
		level: logLevel,
	}),
	new transports.Console({
		level: logLevel,
	}),
];

export const log: Logger = createLogger({
	level: logLevel,
	format: combine(
		errors({ stack: true }),
		timestamp({ format: "YYYY-MM-DD HH:mm:ss.SSS" }),
		logsFormat,
		colorize()
	),
	transports: myTransports,
	rejectionHandlers: [
		new transports.File({ filename: `${logFolder}-rejections.log` }),
	],
});
