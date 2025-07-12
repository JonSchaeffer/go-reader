import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		host: '0.0.0.0',
		port: 3000,
		strictPort: true,
		allowedHosts: ['reader-frontend.porgy-monitor.ts.net'], // Allow all hosts for development
		hmr: {
			clientPort: 3000
		}
	}
});
