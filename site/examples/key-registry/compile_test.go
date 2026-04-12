// Source: site/src/content/docs/components/key-registry.mdx:21,118,134,151
package examples_test

import (
	"testing"

	key "charm.land/bubbles/v2/key"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_KeyRegistryQuickExample verifies the quick example from key-registry.mdx.
func TestCompile_KeyRegistryQuickExample(t *testing.T) {
	registry := teautils.NewKeyRegistry()

	registry.MustRegisterMany([]teautils.KeyMeta{
		{
			ID:        teautils.MustParseKeyIdentifier("nav.up"),
			Binding:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "Move up")),
			HelpModal: true,
			Category:  "Navigation",
			HelpText:  "Move cursor up one line",
		},
		{
			ID:             teautils.MustParseKeyIdentifier("sys.quit"),
			Binding:        key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit")),
			StatusBar:      true,
			StatusBarLabel: "Quit",
			HelpModal:      true,
			Category:       "System",
			HelpText:       "Quit the application",
		},
	})
}

// TestCompile_DisplayFunctions verifies FormatKeyDisplay and GetSortedCategories.
func TestCompile_DisplayFunctions(t *testing.T) {
	km := teautils.KeyMeta{
		Binding: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("Ctrl+C", "copy")),
	}
	label := teautils.FormatKeyDisplay(km)
	_ = label

	registry := teautils.NewKeyRegistry()
	registry.MustRegister(teautils.KeyMeta{
		ID:       teautils.MustParseKeyIdentifier("nav.up"),
		Binding:  key.NewBinding(key.WithKeys("up")),
		Category: "Navigation",
	})

	cats := teautils.GetSortedCategories(
		registry.ByCategory(),
		[]string{"Navigation", "Actions", "System"},
	)
	_ = cats
}

// TestCompile_StatusBarIntegration verifies the status bar integration from key-registry.mdx.
func TestCompile_StatusBarIntegration(t *testing.T) {
	registry := teautils.NewKeyRegistry()
	registry.MustRegister(teautils.KeyMeta{
		ID:             teautils.MustParseKeyIdentifier("sys.quit"),
		Binding:        key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit")),
		StatusBar:      true,
		StatusBarLabel: "Quit",
	})

	statusKeys := registry.ForStatusBar()
	menuItems := make([]teastatus.MenuItem, len(statusKeys))
	for i, km := range statusKeys {
		menuItems[i] = teastatus.NewMenuItem(km.Binding, &teastatus.MenuItemOpts{
			Label: km.StatusBarLabel,
		})
	}
	sb := teastatus.NewStatusBarModel().SetMenuItems(menuItems)
	_ = sb
}

// TestCompile_NamespaceValidation verifies namespace validation from key-registry.mdx.
func TestCompile_NamespaceValidation(t *testing.T) {
	_ = teautils.MustParseKeyIdentifier("nav.up")
	_ = teautils.MustParseKeyIdentifier("action.delete")
	_ = teautils.MustParseKeyIdentifier("sys.help")
}
