import { log } from "../config/logger";
import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { DailyRepository } from "../repository/daily-repo";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { type StockModel, toGoogleFinance } from "../model/googleFinance-model";

export class StockUsecase {
	static async updateDailyStock(user: User): Promise<StockModel[]> {
		try {
			log.info(JSON.stringify(user));

			const stockResponse: Response =
				await GoogleFinanceService.getStocksSummary();
			let stockData: StockModel[] = [];

			if (stockResponse.ok) {
				stockData = toGoogleFinance(await stockResponse.json());
				await DailyRepository.upsert(stockData, user);
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
