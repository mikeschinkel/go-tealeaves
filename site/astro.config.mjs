// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://mikeschinkel.github.io',
	base: '/go-tealeaves',
	integrations: [
		starlight({
			title: 'Tea Leaves',
			logo: { src: './src/assets/logo.png', alt: 'Tea Leaves logo' },
			favicon: '/favicon.png',
			social: [
				{ icon: 'open-book', label: 'Browse Components', href: '/go-tealeaves/components/' },
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/mikeschinkel/go-tealeaves' },
			],
			customCss: ['./src/styles/custom.css'],
			components: {
				Footer: './src/components/Footer.astro',
				PageTitle: './src/components/PageTitle.astro',
			},
			expressiveCode: {
				themes: ['github-dark'],
			},
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Quick Start', slug: 'guides/quick-start' },
						{ label: 'Architecture', slug: 'guides/architecture' },
						{ label: 'Composing Components', slug: 'guides/composition' },
					],
				},
				{
					label: 'Components',
					items: [
						{ label: 'Overview', slug: 'components' },
						{ label: 'teadd — Dropdown', slug: 'components/teadd' },
						{ label: 'teagrid — Data Grid', slug: 'components/teagrid' },
						{ label: 'teamodal — Modals', slug: 'components/teamodal' },
						{ label: 'teanotify — Notifications', slug: 'components/teanotify' },
						{ label: 'teatree — Tree View', slug: 'components/teatree' },
						{ label: 'teatextsel — Text Selection', slug: 'components/teatextsel' },
						{ label: 'teastatus — Status Bar', slug: 'components/teastatus' },
						{ label: 'teadep — Dep Viewer', slug: 'components/teadep' },
						{ label: 'teautils — Utilities', slug: 'components/teautils' },
					],
				},
				{
					label: 'Patterns',
					items: [
						{ label: 'Modal Message Consumption', slug: 'patterns/modal-consumption' },
						{ label: 'Overlay Compositing', slug: 'patterns/overlay-compositing' },
						{ label: 'Wither Methods', slug: 'patterns/wither-methods' },
						{ label: 'Key Registry', slug: 'patterns/key-registry' },
					],
				},
				{
					label: 'Cookbook',
					collapsed: true,
					items: [
						{ label: 'Dropdown in a Form', slug: 'cookbook/dropdown-in-form' },
						{ label: 'Async Notifications', slug: 'cookbook/notification-after-action' },
						{ label: 'Tree + Status Bar', slug: 'cookbook/tree-with-statusbar' },
						{ label: 'Filterable Grid', slug: 'cookbook/grid-with-filtering' },
					],
				},
				{
					label: 'Examples',
					items: [
						{ label: 'Example Gallery', slug: 'examples' },
					],
				},
				{
					label: 'Reference',
					collapsed: true,
					items: [
						{ label: 'Best Practices', slug: 'reference/best-practices' },
						{ label: 'Roadmap', slug: 'reference/roadmap' },
						{ label: 'Charm v2 Migration', slug: 'migration/charm-v2' },
						{ label: 'Contributing', slug: 'contributing' },
					],
				},
			],
		}),
	],
});
