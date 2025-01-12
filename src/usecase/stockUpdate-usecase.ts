import { HTTPException } from "hono/http-exception";
import { type StockModel, toGoogleFinance } from "../model/googleFinance-model";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { DailyRepository } from "../repository/daily-repo";
import type { User } from "@prisma/client";
import { log } from "../config/logger.ts";

export class StockUpdateUsecase {
	static async updateStockData(user: User): Promise<StockModel[]> {
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
