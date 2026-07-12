import { sveltekit } from "@sveltejs/kit/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	server: {
		// In dev the Go backend runs on :8080 (make dev-backend).
		proxy: { "/api": "http://localhost:8080" }
	},
	test: {
		include: ["src/**/*.test.ts"]
	}
});
