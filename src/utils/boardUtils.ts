import { Board } from "@prisma/client";

export class BoardUtils {
	static getBoardName(name: string): Board {
		name = name.toLowerCase();
		if (name.includes("seri b") || name.includes("mvs")) {
			return Board.B_SERIES;
		} else if (name.includes("seri a")) {
			return Board.A_SERIES;
		} else if (name.includes("seri c")) {
			return Board.C_SERIES;
		} else if (name.includes("saham preferen")) {
			return Board.PREFEREN;
		} else {
			return Board.Main;
		}
	}
}
