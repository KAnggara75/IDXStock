export class JsonUtils {
	static safeParseJSON(input: string): unknown {
		try {
			return JSON.parse(input);
		} catch {
			return input;
		}
	}
}
