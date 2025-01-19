import fs from "fs";
import csv from "csv-parser";
import { log } from "../config/logger";
import type { SummaryData } from "../model/summary-model";

export class CsvService {
	static async toJson(request: string): Promise<SummaryData[]> {
		const results: SummaryData[] = [];

		try {
			log.info(request);
			return new Promise((resolve, reject) => {
				fs.createReadStream(request)
					.on("error", (error) => {
						reject(error);
					})
					.pipe(
						csv({
							strict: true,
							separator: ",",
							mapHeaders: ({ header }) =>
								header.toLowerCase().replaceAll(" ", "_"),
						})
					)
					.on("data", (data) => results.push(data))
					.on("end", () => {
						resolve(results);
					});
			});
		} catch {
			log.error("Error");
		}

		return results;
	}
}
