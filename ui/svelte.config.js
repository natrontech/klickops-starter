import adapter from "@sveltejs/adapter-static";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		// Single-page app: the Go backend serves this build output and
		// falls back to index.html for client-side routes.
		adapter: adapter({ fallback: "index.html" })
	}
};

export default config;
