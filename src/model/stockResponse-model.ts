import type { StockModel } from "./googleFinance-model";

export interface StockResponse {
	error?: Error;
	data?: StockModel[];
}

export interface Error {
	code: number;
	message: string;
	status: string;
}
