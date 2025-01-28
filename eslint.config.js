import globals from "globals";
import pluginJs from "@eslint/js";
import tsdoc from "eslint-plugin-tsdoc";
import tseslint from "typescript-eslint";

export default [
	{ ignores: ["**/modules/index.js", "**/src-old/*.{js,mjs,cjs,ts}"] },
	{ files: ["**/src/*.{js,mjs,cjs,ts}"] },
	{
		plugins: {
			tsdoc: tsdoc,
		},
		languageOptions: { globals: globals.node },
		rules: {
			"tsdoc/syntax": "warn",
			"no-unused-vars": "warn",
			"@typescript-eslint/no-unused-vars": "warn",
			"no-console": [
				"error",
				{
					allow: ["warn", "error"],
				},
			],
			semi: [2, "always"],
			indent: ["warn", "tab", { SwitchCase: 1 }],
			quotes: ["warn", "double"],
		},
	},
	pluginJs.configs.recommended,
	...tseslint.configs.recommended,
];
