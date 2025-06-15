import { log } from "../config/logger";
import { type User } from "@prisma/client";
import { prismaClient } from "../config/database";
import type { StockModel } from "../model/googleFinance-model";

export class DailyRepository {
	static async upsert(data: StockModel[], user: User) {
		try {
			const upserts = data.map((stock) =>
				prismaClient.daily.upsert({
					where: {
						code: stock.code,
					},
					update: {
						pe: stock.pe,
						eps: stock.eps,
						low: stock.low,
						high: stock.high,
						low52: stock.low52,
						price: stock.price,
						change: stock.change,
						volume: stock.volume,
						shares: stock.shares,
						high52: stock.high52,
						priceopen: stock.priceopen,
						marketcap: stock.marketcap,
						volumeavg: stock.volumeavg,
						changepct: stock.changepct,
						closeyest: stock.closeyest,
						tradetime: stock.tradetime,
						insertBy: user.username,
					},
					create: {
						pe: stock.pe,
						eps: stock.eps,
						low: stock.low,
						high: stock.high,
						code: stock.code,
						name: stock.name,
						low52: stock.low52,
						price: stock.price,
						change: stock.change,
						volume: stock.volume,
						shares: stock.shares,
						high52: stock.high52,
						priceopen: stock.priceopen,
						marketcap: stock.marketcap,
						volumeavg: stock.volumeavg,
						changepct: stock.changepct,
						closeyest: stock.closeyest,
						tradetime: stock.tradetime,
						insertBy: user.username,
					},
				})
			);

			await Promise.all(upserts);
			log.info(`Daily upsert executed by ${user.username}`);
		} catch (error) {
			log.error(`Error in batch upsert: ${error}`);
		} finally {
			await prismaClient.$disconnect();
		}
	}
}
