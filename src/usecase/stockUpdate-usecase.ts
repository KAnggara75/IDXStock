import { HTTPException } from "hono/http-exception";
import type { StockModel } from "../model/googleFinance-model";
import { GoogleFinanceService } from "../service/googleFinance-service";

export class StockUpdateUsecase {
	static async updateStockData(): Promise<StockModel[]> {
		try {
			const stockData: StockModel[] = await GoogleFinanceService.getStocks();
			return stockData;
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}
}
