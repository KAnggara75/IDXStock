import { log } from "../config/logger";
import {
	toGoogleFinance,
	type GoogleFinance,
	type StockModel,
} from "../model/googleFinance-model";

export class GoogleFinanceService {
	static async getStocks(): Promise<StockModel[]> {
		const apiKey: string = Bun.env.SPREADSHEETS_API_KEY ?? "";
		const spreadsheetId: string = Bun.env.SPREADSHEETS_ID ?? "";
		const sheetName: string = Bun.env.SPREADSHEETS_NAME ?? "Sheet1";

		const url: string | URL =
			`https://sheets.googleapis.com/v4/spreadsheets/${spreadsheetId}/values/${sheetName}?key=${apiKey}`;

		let json: GoogleFinance = {
			values: [],
		};

		let result: StockModel[] = [];

		try {
			const response = await fetch(url);
			json = await response.json();
			result = toGoogleFinance(json);
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (error: any) {
			log.error(error.message);
		}

		return result;
	}
}
