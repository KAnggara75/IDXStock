/* eslint-disable no-console */
import { PrismaClient } from "@prisma/client";
import { StockSeeder, type StockListModel } from "./stock";
const prisma = new PrismaClient();

async function main() {
	const stockList: StockListModel[] = await StockSeeder.getStockList();

	try {
		const upserts = stockList.map((stock) =>
			prisma.stock.upsert({
				where: { code: stock.code },
				update: {
					name: stock.name,
					shares: stock.shares,
					board: stock.board,
					listing_date: stock.listing_date,
				},
				create: {
					code: stock.code,
					name: stock.name,
					shares: stock.shares,
					board: stock.board,
					listing_date: stock.listing_date,
				},
			})
		);

		await Promise.all(upserts); // Jalankan semua upsert secara paralel
		console.info("Batch upsert completed successfully.");
	} catch (error) {
		console.error(error);
		console.error(`Error in batch upsert: ${error}`);
	} finally {
		await prisma.$disconnect();
	}
}

await main();
