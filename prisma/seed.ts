/* eslint-disable no-console */
import { PrismaClient } from "@prisma/client";
import { type StockListModel, StockSeeder } from "./stock";
import uw from "./underwriter.json";

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

		const upsertsUW = uw.map((data) =>
			prisma.underWriter.upsert({
				where: { code: data.code },
				update: {
					name: data.name,
				},
				create: {
					code: data.code,
					name: data.name,
				},
			})
		);

		await Promise.all(upserts); // execute all  upsert in parallel
		console.info("Batch upsert stock completed successfully.");
		await Promise.all(upsertsUW);
		console.info("Batch upsert underWriter completed successfully.");
	} catch (error) {
		console.error(error);
		console.error(`Error in batch upsert: ${error}`);
	} finally {
		await prisma.$disconnect();
	}
}

await main();
