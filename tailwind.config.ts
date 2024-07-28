import type { Config } from 'tailwindcss';

export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],

	theme: {
		extend: {
			colors: {
				"primary": "#0B032D",
				"secondary": "#CB974A",
				"highlight": "#BD1D1F",
				"background": "#EAEAEA"
			}
		}
	},

	plugins: []
} as Config;
