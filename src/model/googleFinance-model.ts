import { HTTPException } from "hono/http-exception";
import { AppConstant } from "../config/appConstant";
import { NumberUtil } from "../utils/numberUtils";
import { DateUtils } from "../utils/dateUtils";

export interface GoogleFinance {
	range?: string;
	majorDimension?: string;
	error?: Error;
	values?: string[][];
}

export interface Error {
	code: number;
	message: string;
	status: string;
}

export interface StockModel {
	code: string; // "AALI",
	name: string; // "Astra Agro Lestari Tbk PT",
	price: number | null; // "5925",
	priceopen: number | null; // "5975",
	high: number | null; // "5975",
	low: number | null; // "5900",
	volume: number | null; // "267400",
	marketcap: number | null; // "11403770475000",
	tradetime: Date; // "2025-01-10 16:09:51",
	volumeavg: number | null; // "383427",
	pe: number | null; // "10,79",
	eps: number | null; // "548,90",
	high52: number | null; // "7200",
	low52: number | null; // "5250",
	change: number | null; // "-50",
	changepct: number | null; // "-0,84",
	closeyest: number | null; // "5975",
	shares: number | null; // "1924688333"
}

export function toGoogleFinance(data: GoogleFinance): StockModel[] {
	const value: string[][] | undefined = data.values;
	const result: StockModel[] = [];

	if (value) {
		const header = value[0].map((v) => v.toLowerCase());

		if (header.length != AppConstant.GOOGLE_FINANCE_FORMAT.length) {
			throw new HTTPException(422, {
				message: "Bad data GOOGLE_FINANCE_FORMAT",
			});
		}

		for (let i = 0; header.length > i; i++) {
			if (header[i] !== AppConstant.GOOGLE_FINANCE_FORMAT[i]) {
				throw new HTTPException(422, {
					message: "Bad data GOOGLE_FINANCE_FORMAT",
				});
			}
		}

		for (let i = 1; value.length > i; i++) {
			const stock = value[i].map((v) => v.replaceAll("#N/A", ""));
			const stockData: StockModel = {
				code: stock[0],
				name: stock[1] === "" ? stock[0] : stock[1],
				price: NumberUtil.toInt(stock[2]),
				priceopen: NumberUtil.toInt(stock[3]),
				high: NumberUtil.toInt(stock[4]),
				low: NumberUtil.toInt(stock[5]),
				volume: NumberUtil.toInt(stock[6]),
				marketcap: NumberUtil.toInt(stock[7]),
				tradetime: DateUtils.toDate(stock[8]),
				volumeavg: NumberUtil.toInt(stock[9]),
				pe: NumberUtil.toFloat(stock[10]),
				eps: NumberUtil.toFloat(stock[11]),
				high52: NumberUtil.toInt(stock[12]),
				low52: NumberUtil.toInt(stock[13]),
				change: NumberUtil.toInt(stock[14]),
				changepct: NumberUtil.toFloat(stock[15]),
				closeyest: NumberUtil.toInt(stock[16]),
				shares: NumberUtil.toInt(stock[17]),
			};

			result.push(stockData);
		}
	} else {
		throw new HTTPException(422, {
			message: "Bad data GoogleFinance",
		});
	}

	return result;
}
