import { createLogger, format, transports } from "winston";

const { combine, timestamp, printf, colorize } = format;

const logsFormat = printf(({ level, message, timestamp }) => {
	return `${timestamp}|[${level.toUpperCase()}]|${message}|`;
});

const logLevel: string = Bun.env.LOG_LEVEL ?? "warn";
const todayDate: string = new Date().toISOString().slice(0, 10);
const logFolder: string = `./logs/${todayDate.replace(/-/g, "")}/${todayDate}`;

const myTransports: (
	| transports.ConsoleTransportInstance
	| transports.FileTransportInstance
)[] = [
	new transports.File({
		filename: `${logFolder}-error.log`,
		level: "error",
	}),
	new transports.File({
		filename: `${logFolder}-warn.log`,
		level: "warn",
	}),
];

if (logLevel == "warn") {
	myTransports.push(
		new transports.Console({
			level: "warn",
		})
	);
} else if (logLevel == "info") {
	myTransports.push(
		new transports.Console({
			level: "info",
		})
	);
} else {
	myTransports.push(
		new transports.Console({
			level: "error",
		})
	);
}

export const log = createLogger({
	format: combine(
		timestamp({ format: "YYYY-MM-DD HH:mm:ss.SSS" }),
		logsFormat,
		colorize()
	),
	transports: myTransports,
	rejectionHandlers: [
		new transports.File({ filename: `${logFolder}-rejections.log` }),
	],
});
