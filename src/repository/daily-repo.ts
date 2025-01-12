import { PrismaClient, type User } from "@prisma/client";
import type { StockModel } from "../model/googleFinance-model";
import { log } from "../config/logger";

const prisma = new PrismaClient();

export class DailyRepository {
	static async upsert(data: StockModel[], user: User) {
		try {
			const upserts = data.map((stock) =>
				prisma.daily.upsert({
					where: { code: stock.code },
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

			await Promise.all(upserts); // Jalankan semua upsert secara paralel
			log.info("Batch upsert completed successfully.");
		} catch (error) {
			log.error(`Error in batch upsert: ${error}`);
		} finally {
			await prisma.$disconnect();
		}
	}
}
