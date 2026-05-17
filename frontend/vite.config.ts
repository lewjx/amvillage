import { svelte } from "@sveltejs/vite-plugin-svelte"
import { defineConfig } from "vite"
import postcss from "./postcss.config.js"

export default defineConfig({
	plugins: [svelte()],
	css: {
		postcss,
	},
	base: "",
	server: {
		proxy: {
			"/ws": {
				target: "http://127.0.0.1:8080",
				ws: true,
			},
		},
	},
})
