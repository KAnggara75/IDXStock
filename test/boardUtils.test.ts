import { describe, expect, it } from "bun:test";
import { Board } from "@prisma/client";
import { BoardUtils } from "../src/utils/boardUtils";

describe("BoardUtils Test", () => {
	it.each([
		["Saham MVS", Board.B_SERIES],
		["Saham Main", Board.Main],
		["Saham Seri A", Board.A_SERIES],
		["Saham Seri B", Board.B_SERIES],
		["Saham Seri C", Board.C_SERIES],
		["Saham Preferen", Board.PREFEREN],
		["SAHAM Preferen", Board.PREFEREN],
	])("%s => %s", (name: string, expected: Board) => {
		const board: Board = BoardUtils.getBoardName(name);
		expect(board).toBe(expected);
	});
});
