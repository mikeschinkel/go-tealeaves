// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// Package-to-slug mapping (canonical reference for evaluator verification):
// teacolor      -> components/color-constants         [foundation]
// teacrumbs     -> components/breadcrumb-nav
// teadiff       -> components/diff-viewer
// teafields     -> components/dropdown-control
// teagrid       -> components/grid-view
// teaguide      -> components/guide-overlay
// teahelp       -> components/help-visor
// teahilite     -> components/syntax-highlighting     [foundation]
// tealayout     -> components/layout-engine
// teamodal      -> components/choice-dialog, components/confirm-dialog, components/list-dialog, components/multiselect-dialog, components/progress-dialog
// teanotify     -> components/notification-view
// teapane       -> components/pane-widgets            [foundation]
// teastatus     -> components/statusbar-view
// teatext       -> components/text-selection
// teatree       -> components/tree-view, components/drilldown-view
// teautils      -> components/key-registry, components/positioning, components/theming  [foundation]

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
				{ icon: 'open-book', label: 'Documentation', href: '/go-tealeaves/guides/quick-start/' },
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
						{ label: 'Modal Message Consumption', slug: 'guides/modal-consumption' },
						{ label: 'Wither Methods', slug: 'guides/wither-methods' },
					],
				},
				{
					label: 'Components',
					items: [
						{ label: 'Overview', slug: 'components' },
						{
							label: 'Views',
							items: [
								{ label: 'Grid View', slug: 'components/grid-view' },
								{ label: 'Tree View', slug: 'components/tree-view' },
								{ label: 'Drilldown View', slug: 'components/drilldown-view' },
								{ label: 'Status Bar', slug: 'components/statusbar-view' },
								{ label: 'Notifications', slug: 'components/notification-view' },
								{ label: 'Diff Viewer', slug: 'components/diff-viewer' },
								{ label: 'Breadcrumb Nav', slug: 'components/breadcrumb-nav' },
							],
						},
						{
							label: 'Dialogs',
							items: [
								{ label: 'Confirm Dialog', slug: 'components/confirm-dialog' },
								{ label: 'Choice Dialog', slug: 'components/choice-dialog' },
								{ label: 'List Dialog', slug: 'components/list-dialog' },
								{ label: 'Progress Dialog', slug: 'components/progress-dialog' },
								{ label: 'MultiSelect Dialog', slug: 'components/multiselect-dialog' },
								{ label: 'Guide Overlay', slug: 'components/guide-overlay' },
							],
						},
						{
							label: 'Controls',
							items: [
								{ label: 'Dropdown Control', slug: 'components/dropdown-control' },
							],
						},
						{
							label: 'Text',
							items: [
								{ label: 'Text Selection', slug: 'components/text-selection' },
							],
						},
						{
							label: 'Layout',
							items: [
								{ label: 'Layout Engine', slug: 'components/layout-engine' },
							],
						},
						{
							label: 'System',
							items: [
								{ label: 'Help Visor', slug: 'components/help-visor' },
								{ label: 'Pane Widgets', slug: 'components/pane-widgets' },
								{ label: 'Key Registry', slug: 'components/key-registry' },
								{ label: 'Theming', slug: 'components/theming' },
								{ label: 'Color Constants', slug: 'components/color-constants' },
								{ label: 'Syntax Highlighting', slug: 'components/syntax-highlighting' },
								{ label: 'Positioning', slug: 'components/positioning' },
							],
						},
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
						{ label: 'Module Reference', slug: 'reference/modules' },
						{ label: 'Roadmap', slug: 'reference/roadmap' },
						{ label: 'Overlay Compositing', slug: 'reference/overlay-compositing' },
						{ label: 'Charm v2 Migration', slug: 'migration/charm-v2' },
							{ label: 'Contributing', slug: 'contributing' },
						{ label: 'Why Tea Leaves?', slug: 'guides/why-tea-leaves' },
					],
				},
			],
		}),
	],
});
