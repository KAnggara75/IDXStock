export interface StockData {
	range?: string;
	majorDimension?: string;
	values: StockModel[][];
	error?: StockError;
}

export interface StockError {
	code: number;
	message: string;
	status: string;
	details: Detail[];
}

interface Detail {
	"@type": string;
	reason?: string;
	domain?: string;
	locale?: string;
	message?: string;
}

interface StockModel {
	CODE: string;
	name: string;
	price: string;
	priceopen: string;
	high: string;
	low: string;
	volume: string;
	marketcap: string;
	tradetime: string;
	volumeavg: string;
	pe: string;
	eps: string;
	high52: string;
	low52: string;
	change: string;
	changepct: string;
	closeyest: string;
	shares: string;
}
