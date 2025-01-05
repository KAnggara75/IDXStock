import globals from "globals";
import pluginJs from "@eslint/js";
import tseslint from "typescript-eslint";

/** @type {import('eslint').Linter.Config[]} */
export default [
	{ ignores: ["**/modules/index.js", "**/src-old/*.{js,mjs,cjs,ts}"] },
	{ files: ["**/src/*.{js,mjs,cjs,ts}"] },
	{
		languageOptions: { globals: globals.node },
		rules: {
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
