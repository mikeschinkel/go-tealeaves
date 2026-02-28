// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://mikeschinkel.github.io',
	base: '/go-tealeaves',
	integrations: [
		starlight({
			title: 'go-tealeaves',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/mikeschinkel/go-tealeaves' }],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Installation', slug: 'guides/getting-started' },
						{ label: 'Architecture', slug: 'guides/architecture' },
					],
				},
				{
					label: 'Components',
					autogenerate: { directory: 'reference' },
				},
				{
					label: 'Migration',
					items: [
						{ label: 'Charm v2 Migration', slug: 'guides/charm-v2-migration' },
					],
				},
			],
		}),
	],
});
