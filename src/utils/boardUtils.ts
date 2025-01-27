import { Board } from "@prisma/client";

export class BoardUtils {
	static getBoardName(name: string): Board {
		if (name.includes("Seri B") || name.includes("MVS")) {
			return Board.B_SERIES;
		} else if (name.includes("Seri A")) {
			return Board.A_SERIES;
		} else if (name.includes("Seri C")) {
			return Board.C_SERIES;
		} else if (name.includes("Saham Preferen")) {
			return Board.PREFEREN;
		} else {
			return Board.Main;
		}
	}
}
