import { describe, it, expect, spyOn } from "bun:test";
import { log } from "../src/config/logger";

describe("Logger Test", () => {
	it("Log error Test", async () => {
		const logMsg = "This is an error log Test";
		const errorLog = spyOn(log, "error");
		expect(errorLog).toHaveBeenCalledTimes(0);
		log.error(logMsg);
		expect(errorLog).toHaveBeenCalledTimes(1);
		expect(errorLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log warn Test", async () => {
		const logMsg = "This is an warn log Test";
		const warnLog = spyOn(log, "warn");
		expect(warnLog).toHaveBeenCalledTimes(0);
		log.warn(logMsg);
		expect(warnLog).toHaveBeenCalledTimes(1);
		expect(warnLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log info Test", async () => {
		const logMsg = "This is an info log Test";
		const infoLog = spyOn(log, "info");
		expect(infoLog).toHaveBeenCalledTimes(0);
		log.info(logMsg);
		expect(infoLog).toHaveBeenCalledTimes(1);
		expect(infoLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log debug Test", async () => {
		const logMsg = "This is an debug log Test";
		const debugLog = spyOn(log, "debug");
		expect(debugLog).toHaveBeenCalledTimes(0);
		log.debug(logMsg);
		expect(debugLog).toHaveBeenCalledTimes(1);
		expect(debugLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log http Test", async () => {
		const logMsg = "This is an http log Test";
		const httpLog = spyOn(log, "http");
		expect(httpLog).toHaveBeenCalledTimes(0);
		log.http(logMsg);
		expect(httpLog).toHaveBeenCalledTimes(1);
		expect(httpLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log verbose Test", async () => {
		const logMsg = "This is an verbose log Test";
		const verboseLog = spyOn(log, "verbose");
		expect(verboseLog).toHaveBeenCalledTimes(0);
		log.verbose(logMsg);
		expect(verboseLog).toHaveBeenCalledTimes(1);
		expect(verboseLog.mock.calls).toEqual([[logMsg]]);
	});

	it("Log silly Test", async () => {
		const logMsg = "This is an silly log Test";
		const sillyLog = spyOn(log, "silly");
		expect(sillyLog).toHaveBeenCalledTimes(0);
		log.silly(logMsg);
		expect(sillyLog).toHaveBeenCalledTimes(1);
		expect(sillyLog.mock.calls).toEqual([[logMsg]]);
	});
});
