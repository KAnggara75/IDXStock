import { Board, type Stock } from "@prisma/client";

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

		const result: Stock[] = [];

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
					const stockData: Stock = {
						code: stock[1],
						name: stock[2] === "" ? stock[1] : stock[2],
						listing_date: StockSeeder.formatDate(stock[3]),
						delisting_date: null,
						shares: BigInt(stock[4].replaceAll(",", "")),
						board: StockSeeder.toBoardEnum(stock[5]),
					};

					result.push(stockData);
				}
			}
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (error: any) {
			console.error(error.message);
		}

		const finn: Stock = {
			board: Board.Main,
			code: "FINN",
			listing_date: StockSeeder.formatDate("29 Jul 2019"),
			delisting_date: StockSeeder.formatDate("02 Mar 2021"),
			name: "First Indo American Leasing Tbk",
			shares: BigInt(2188500000),
		};
		result.push(finn);

		const cntb: Stock = {
			board: Board.Main,
			code: "CNTB",
			listing_date: StockSeeder.formatDate("12 Jul 2021"),
			delisting_date: null,
			name: "Saham Seri B ( Centex Tbk )",
			shares: BigInt(2188500000),
		};
		result.push(cntb);

		return result;
	}

	static formatDate(inputDate: string) {
		return new Date(
			new Date(inputDate.replace("Agt", "Aug")).getTime() + 7 * 60 * 60 * 1000
		);
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
