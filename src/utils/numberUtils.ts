export class NumberUtil {
	/**
	 * Convert {@link number} or {@link string} to {@link integer}
	 *
	 * @param number - A {@link number} or {@link string} expression.
	 *
	 * @returns {@link number}
	 *
	 * @example
	 * Here's a {@link number} example:
	 * ```
	 * NumberUtil.toInt(1); // result 1
	 * NumberUtil.toInt(2); // result 2
	 * NumberUtil.toInt(3.1); // result 3
	 * NumberUtil.toInt(4.9); // result 4
	 * NumberUtil.toInt(0.9); // result 0
	 * ```
	 *
	 * @example
	 * Here's a {@link string} example:
	 * ```
	 * NumberUtil.toInt(""); // result 0
	 * NumberUtil.toInt("2"); // result 2
	 * NumberUtil.toInt("a"); // result 0
	 * NumberUtil.toInt("0,1"); // result 0
	 * NumberUtil.toInt("1,2"); // result 1
	 * NumberUtil.toInt("2,1"); // result 2
	 * NumberUtil.toInt("1m3"); // result 0
	 * NumberUtil.toInt("1.15503E+11"); // result 115503000000
	 * NumberUtil.toInt("1,15503E+11"); // result 115503000000
	 * ```
	 */
	static toInt(number: string | number): number {
		let result: number = 0;
		try {
			if (typeof number === "number") {
				return Number(Math.floor(number));
			} else if (number.includes("E+") || number.includes(",")) {
				result = Number(Math.floor(Number(number.replaceAll(",", "."))));
			} else if (number === "" || isNaN(Number(number))) {
				result = 0;
			} else {
				result = parseInt(number.replaceAll(",", "."));
			}
			return isNaN(result) ? 0 : result;
		} catch {
			return 0;
		}
	}

	static toFloat(number: string | number): number {
		let result: number = 0;

		try {
			if (typeof number === "number") {
				return number;
			} else if (number.includes("E+") || number.includes(",")) {
				result = Number(Number(number.replaceAll(",", ".")));
			} else if (number === "" || isNaN(Number(number))) {
				result = 0.0;
			} else {
				result = parseFloat(number.replaceAll(",", "."));
			}
			return isNaN(result) ? 0 : result;
		} catch {
			return 0.0;
		}
	}

	/**
	 * Convert {@link number} or {@link string} to {@link BigInt}
	 *
	 * @param number - A {@link number} or {@link string} expression.
	 *
	 * @returns {@link BigInt}
	 *
	 * @example
	 * Here's a {@link number} example:
	 * ```
	 * NumberUtil.toBigInt(1); // result 1
	 * NumberUtil.toBigInt(2); // result 2
	 * NumberUtil.toBigInt(3.1); // result 3
	 * NumberUtil.toBigInt(4.9); // result 4
	 * NumberUtil.toBigInt(0.9); // result 0
	 * ```
	 *
	 * @example
	 * Here's a {@link string} example:
	 * ```
	 * NumberUtil.toBigInt(""); // result 0
	 * NumberUtil.toBigInt("2"); // result 2
	 * NumberUtil.toBigInt("a"); // result 0
	 * NumberUtil.toBigInt("0,1"); // result 0
	 * NumberUtil.toBigInt("1,2"); // result 1
	 * NumberUtil.toBigInt("2,1"); // result 2
	 * NumberUtil.toBigInt("1m3"); // result 0
	 * NumberUtil.toBigInt("1.15503E+11"); // result 115503000000
	 * NumberUtil.toBigInt("1,15503E+11"); // result 115503000000
	 * ```
	 */
	static toBigInt(number: string | number): bigint {
		try {
			if (typeof number === "number") {
				return BigInt(Math.floor(number));
			} else if (number.includes("E+") || number.includes(",")) {
				return BigInt(Math.floor(Number(number.replaceAll(",", "."))));
			} else if (number === "") {
				return BigInt(0);
			} else {
				return BigInt(Math.floor(Number(number.replaceAll(",", "."))));
			}
		} catch {
			return BigInt(0);
		}
	}
}
