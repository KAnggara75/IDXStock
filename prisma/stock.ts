import { Board } from "@prisma/client";

export class StockSeeder {
	static async getStockList() {
		const apiKey: string = Bun.env.SPREADSHEETS_API_KEY ?? "";
		const spreadsheetId: string = Bun.env.SPREADSHEETS_ID ?? "";
		const stockListKey: string[] = [
			"No",
			"Code",
			"Company Name",
			"Listing Date",
			"Shares",
			"Listing Board",
		];

		const url: string | URL =
			`https://sheets.googleapis.com/v4/spreadsheets/${spreadsheetId}/values/List?key=${apiKey}`;

		const result: StockListModel[] = [];

		try {
			const response: Response = await fetch(url);
			const json: GoogleFinance = await response.json();
			const stockList: string[][] | undefined = json.values;

			if (!response.ok) {
				throw new Error(json.error?.message);
			}

			if (stockList) {
				const header = stockList[0];

				if (header.length != stockListKey.length) {
					throw new Error("Bad data GOOGLE_FINANCE_FORMAT");
				}

				for (let i = 0; header.length > i; i++) {
					if (header[i] !== stockListKey[i]) {
						throw new Error("Bad data GOOGLE_FINANCE_FORMAT");
					}
				}

				for (let i = 1; stockList.length > i; i++) {
					const stock = json.values![i].map((v) => v.replaceAll("#N/A", ""));
					const stockData: StockListModel = {
						code: stock[1],
						name: stock[2] === "" ? stock[1] : stock[2],
						listing_date: StockSeeder.formatDate(stock[3]),
						shares: parseInt(stock[4].replaceAll(",", "")),
						board: StockSeeder.toBoardEnum(stock[5]),
					};

					result.push(stockData);
				}
			}
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (error: any) {
			console.error(error.message);
		}

		return result;
	}

	static formatDate(inputDate: string) {
		const date = new Date(inputDate.replaceAll("Agt", "Aug"));
		return date.toISOString();
	}

	static toBoardEnum(board: string): Board {
		switch (board.replaceAll(" ", "")) {
			case "Main":
				return Board.Main;
			case "Watchlist":
				return Board.Watchlist;
			case "Development":
				return Board.Development;
			case "Acceleration":
				return Board.Acceleration;
			case "EkonomiBaru":
				return Board.EkonomiBaru;
			default:
				return Board.Main;
		}
	}
}

interface GoogleFinance {
	range?: string;
	majorDimension?: string;
	error?: Error;
	values?: string[][];
}

export interface StockListModel {
	code: string; // "AALI",
	name: string; // "Astra Agro Lestari Tbk PT",
	listing_date: string; // "09 Dec 1997"
	shares: number; // 1924688333
	board: Board; // "Main",
}
