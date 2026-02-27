// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://mikeschinkel.github.io',
	base: '/go-tealeaves',
	integrations: [
		starlight({
			title: 'Tea Leaves for Go + Bubble Tea',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/mikeschinkel/go-tealeaves' }],
			sidebar: [
				{
					label: 'Guides',
					items: [
						// Each item here is one entry in the navigation menu.
						{ label: 'Example Guide', slug: 'guides/example' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
