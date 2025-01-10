import { HTTPException } from "hono/http-exception";
import { log } from "../config/logger";

export class GoogleFinanceService {
	static async getStocks(): Promise<Response> {
		const apiKey: string = Bun.env.SPREADSHEETS_API_KEY ?? "";
		const spreadsheetId: string = Bun.env.SPREADSHEETS_ID ?? "";
		const sheetName: string = Bun.env.SPREADSHEETS_NAME ?? "Sheet1";

		const url: string | URL =
			`https://sheets.googleapis.com/v4/spreadsheets/${spreadsheetId}/values/${sheetName}?key=${apiKey}`;

		try {
			const response: Response = await fetch(url);
			return response;
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (error: any) {
			log.error(error.message);
			throw new HTTPException(500, {
				message: error.message,
			});
		}
	}
}
