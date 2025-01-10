import { HTTPException } from "hono/http-exception";
import { toGoogleFinance, type StockModel } from "../model/googleFinance-model";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { DailyRepository } from "../repository/daily-repo";

export class StockUpdateUsecase {
	static async updateStockData(): Promise<StockModel[]> {
		
		try {
			const stockResponse: Response = await GoogleFinanceService.getStocks();
			let stockData: StockModel[] = [];

			if (stockResponse.ok) {
				stockData = toGoogleFinance(await stockResponse.json());
				await DailyRepository.upsert(stockData);
				return stockData;
			} else {
				throw new HTTPException(400, {
					message: "bad request file name",
				});
			}
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}
}
