import { log } from "../config/logger.ts";
import { DateUtils } from "../utils/dateUtils.ts";
import { prismaClient } from "../config/database.ts";
import type { SummaryData } from "../model/summary-model.ts";
import { NumberUtil } from "../utils/numberUtils.ts";
import { BoardUtils } from "../utils/boardUtils.ts";

export class HistoryRepository {
	static async insertSummary(data: SummaryData[]) {
		try {
			for (const stock of data) {
				await prismaClient.stock
					.upsert({
						where: {
							code: stock.stock_code,
						},
						update: {
							name: stock.company_name,
						},
						create: {
							code: stock.stock_code,
							name: stock.company_name,
							shares: stock.listed_shares,
							board: BoardUtils.getBoardName(stock.company_name),
						},
					})
					.catch((error) => {
						log.error("Error upsert data:", error);
					});

				const stockData = await prismaClient.history
					.findFirst({
						where: {
							code: stock.stock_code,
							date: DateUtils.toDate2(stock.last_trading_date),
						},
					})
					.catch((error) => {
						log.error("Error Find data:", error);
					});

				if (!stockData) {
					log.debug(
						`Insert CODE ${stock.stock_code}-${stock.last_trading_date}`
					);
					await prismaClient.history
						.create({
							data: {
								code: stock.stock_code,
								date: DateUtils.toDate2(stock.last_trading_date),
								previous: NumberUtil.toInt(stock.previous),
								open_price: NumberUtil.toInt(stock.open_price),
								first_trade: NumberUtil.toInt(stock.first_trade),
								high: NumberUtil.toInt(stock.high),
								low: NumberUtil.toInt(stock.low),
								close: NumberUtil.toInt(stock.close),
								change: NumberUtil.toInt(stock.change),
								volume: NumberUtil.toInt(stock.volume),
								value: NumberUtil.toBigInt(stock.value),
								frequency: NumberUtil.toBigInt(stock.frequency),
								index_individual: NumberUtil.toFloat(stock.index_individual),
								offer: NumberUtil.toBigInt(stock.offer),
								offer_volume: NumberUtil.toBigInt(stock.offer_volume),
								bid: NumberUtil.toBigInt(stock.bid),
								bid_volume: NumberUtil.toBigInt(stock.bid_volume),
								listed_shares: NumberUtil.toBigInt(stock.listed_shares),
								tradeble_shares: NumberUtil.toBigInt(stock.tradeble_shares),
								weight_for_index: NumberUtil.toBigInt(stock.weight_for_index),
								foreign_sell: NumberUtil.toBigInt(stock.foreign_sell),
								foreign_buy: NumberUtil.toBigInt(stock.foreign_buy),
								non_regular_volume: NumberUtil.toBigInt(
									stock.non_regular_volume
								),
								non_regular_value: NumberUtil.toBigInt(stock.non_regular_value),
								non_regular_frequency: NumberUtil.toBigInt(
									stock.non_regular_frequency
								),
							},
						})
						.catch((error) => {
							log.error("Error inserting data:", error);
						});
				}
			}
		} catch (error) {
			log.error(`Error in insert stock summary: ${error}`);
		} finally {
			await prismaClient.$disconnect();
		}
	}
}
