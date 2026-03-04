package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/teadrpdwn"
)

const (
	modeConfiguring   = "configuring"
	modeDemonstrating = "demonstrating"
)

type exampleModel struct {
	// Current mode
	mode string

	// Configuration state
	currentField   int
	configComplete [5]bool // Track which fields are configured (now 5 fields)

	// Selector options
	horizontalOptions []string
	horizontalValues  []int // 0=Left, 1=Middle, 2=Right
	horizontalIdx     int

	verticalOptions []string
	verticalValues  []int // 0=Top, 1=Middle, 2=Bottom
	verticalIdx     int

	numItemsOptions []string
	numItemsValues  []int
	numItemsIdx     int

	widthOptions []string
	widthValues  []string
	widthIdx     int

	startOpenOptions []string
	startOpenValues  []bool
	startOpenIdx     int

	// Selected configuration values
	selectedHorizontal int // 0=Left, 1=Middle, 2=Right
	selectedVertical   int // 0=Top, 1=Middle, 2=Bottom
	selectedNumItems   int
	selectedWidth      string
	selectedStartOpen  bool

	// Demo state
	dropdown      teadrpdwn.DropdownModel
	selectedValue string // Currently selected dropdown value
	fieldRow      int    // Field position (independent of dropdown position)
	fieldCol      int    // Field column position
	screenWidth   int
	screenHeight  int
}

func main() {
	// Ensure that term.GetSize() is initialized before continuing.
	// This is needed in GoLand terminal for debugging, but is not harmful if not needed.
	teadrpdwn.EnsureTermGetSize(os.Stdout.Fd())

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		cliutil.Stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}

func initialModel() exampleModel {
	return exampleModel{
		mode:         modeConfiguring,
		currentField: 0,

		// Horizontal position options
		horizontalOptions: []string{"Left", "Middle", "Right"},
		horizontalValues:  []int{0, 1, 2},

		// Vertical position options
		verticalOptions: []string{"Top", "Middle", "Bottom"},
		verticalValues:  []int{0, 1, 2},

		// Number of items options
		numItemsOptions: []string{"3 items", "7 items", "25 items"},
		numItemsValues:  []int{3, 7, 25},

		// Width options (will show actual character counts)
		widthOptions: []string{"Short (10 chars)", "Medium (calculated)", "Long (calculated)"},
		widthValues:  []string{"short", "medium", "long"},

		// Start open options
		startOpenOptions: []string{"No", "Yes"},
		startOpenValues:  []bool{false, true},
	}
}

func (m exampleModel) Init() tea.Cmd {
	return nil
}

func (m exampleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == modeConfiguring {
		return m.updateConfiguring(msg)
	}
	return m.updateDemonstrating(msg)
}

func (m exampleModel) updateConfiguring(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			// Go back to previous field
			if m.currentField > 0 {
				m.currentField--
				m.configComplete[m.currentField] = false
			}
			return m, nil

		case "up", "k":
			// Navigate up in current selector
			switch m.currentField {
			case 0:
				m.horizontalIdx--
				if m.horizontalIdx < 0 {
					m.horizontalIdx = len(m.horizontalOptions) - 1
				}
			case 1:
				m.verticalIdx--
				if m.verticalIdx < 0 {
					m.verticalIdx = len(m.verticalOptions) - 1
				}
			case 2:
				m.numItemsIdx--
				if m.numItemsIdx < 0 {
					m.numItemsIdx = len(m.numItemsOptions) - 1
				}
			case 3:
				m.widthIdx--
				if m.widthIdx < 0 {
					m.widthIdx = len(m.widthOptions) - 1
				}
			case 4:
				m.startOpenIdx--
				if m.startOpenIdx < 0 {
					m.startOpenIdx = len(m.startOpenOptions) - 1
				}
			}
			return m, nil

		case "down", "j":
			// Navigate down in current selector
			switch m.currentField {
			case 0:
				m.horizontalIdx++
				if m.horizontalIdx >= len(m.horizontalOptions) {
					m.horizontalIdx = 0
				}
			case 1:
				m.verticalIdx++
				if m.verticalIdx >= len(m.verticalOptions) {
					m.verticalIdx = 0
				}
			case 2:
				m.numItemsIdx++
				if m.numItemsIdx >= len(m.numItemsOptions) {
					m.numItemsIdx = 0
				}
			case 3:
				m.widthIdx++
				if m.widthIdx >= len(m.widthOptions) {
					m.widthIdx = 0
				}
			case 4:
				m.startOpenIdx++
				if m.startOpenIdx >= len(m.startOpenOptions) {
					m.startOpenIdx = 0
				}
			}
			return m, nil

		case "enter":
			// Confirm current selection and move to next
			switch m.currentField {
			case 0:
				m.selectedHorizontal = m.horizontalValues[m.horizontalIdx]
				m.configComplete[0] = true
				m.currentField++
			case 1:
				m.selectedVertical = m.verticalValues[m.verticalIdx]
				m.configComplete[1] = true
				m.currentField++
			case 2:
				m.selectedNumItems = m.numItemsValues[m.numItemsIdx]
				m.configComplete[2] = true
				m.currentField++
			case 3:
				m.selectedWidth = m.widthValues[m.widthIdx]
				m.configComplete[3] = true
				// Update width options with actual character counts now that we have screen size
				m = m.updateWidthOptions()
				m.currentField++
			case 4:
				m.selectedStartOpen = m.startOpenValues[m.startOpenIdx]
				m.configComplete[4] = true
				// All configured - switch to demo mode
				m.mode = modeDemonstrating
				m = m.setupDemo()
				return m, nil
			}
			return m, nil
		}
	}

	return m, nil
}

func (m exampleModel) updateDemonstrating(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var dropdown tea.Model

	// Handle dropdown messages first
	switch msg := msg.(type) {
	case teadrpdwn.OptionSelectedMsg:
		// Update selected value
		m.selectedValue = msg.Text
		m.dropdown, cmd = m.dropdown.Close()
		return m, cmd
	}

	// Let dropdown handle input
	dropdown, cmd = m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teadrpdwn.DropdownModel)
		return m, cmd
	}

	// Dropdown didn't handle - parent processes
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Return to configuration (only when dropdown is closed)
			if !m.dropdown.IsOpen {
				m.mode = modeConfiguring
				m.currentField = 0
				m.configComplete = [5]bool{false, false, false, false, false}
				return m, nil
			}
		case "r":
			// Return to configuration
			m.mode = modeConfiguring
			m.currentField = 0
			m.configComplete = [5]bool{false, false, false, false, false}
			return m, nil
		case "space", "enter", "down", "right":
			// Open dropdown if closed
			if !m.dropdown.IsOpen {
				m.dropdown, cmd = m.dropdown.Open()
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		m.dropdown = m.dropdown.WithScreenSize(msg.Width, msg.Height)
		// Recalculate field position
		m = m.setupDemo()
		return m, nil
	}

	return m, nil
}

// updateWidthOptions updates width options with actual character counts
func (m exampleModel) updateWidthOptions() exampleModel {
	if m.screenWidth > 0 {
		mediumChars := (m.screenWidth * 70) / 100
		longChars := (m.screenWidth * 125) / 100
		m.widthOptions = []string{
			"Short (10 chars)",
			fmt.Sprintf("Medium (%d chars)", mediumChars),
			fmt.Sprintf("Long (%d chars)", longChars),
		}
	}
	return m
}

// getWidthDescription returns the current width description
func (m exampleModel) getWidthDescription() string {
	if m.screenWidth > 0 {
		switch m.selectedWidth {
		case "short":
			return "10 chars"
		case "medium":
			return fmt.Sprintf("%d chars", (m.screenWidth*70)/100)
		case "long":
			return fmt.Sprintf("%d chars", (m.screenWidth*125)/100)
		}
	}
	return m.selectedWidth
}

func (m exampleModel) setupDemo() exampleModel {
	var cmd tea.Cmd

	// Update width options with actual character counts
	m = m.updateWidthOptions()

	// Generate items based on configuration
	items := generateItems(m.selectedWidth, m.selectedNumItems, m.screenWidth)

	// Initialize with first item selected
	if len(items) > 0 {
		m.selectedValue = items[0].Text
	}

	// Calculate base field position (in screen coordinates)
	fieldRow, fieldCol := calculateFieldPosition(m.selectedHorizontal, m.selectedVertical, m.screenWidth, m.screenHeight)

	// For middle positions, adjust field column to center the text
	// Work in interior coordinates for easier calculation
	triangle := "▶"
	fieldText := triangle + " " + m.selectedValue
	fieldTextWidth := ansi.StringWidth(fieldText)
	interiorWidth := m.screenWidth - 2

	// Convert fieldCol from screen coords to interior coords
	interiorFieldCol := fieldCol - 1 // -1 for left border

	adjustedFieldCol := interiorFieldCol
	if m.selectedHorizontal == 1 { // Middle
		adjustedFieldCol = interiorFieldCol - fieldTextWidth/2
	}

	// Ensure field doesn't extend past right edge
	if adjustedFieldCol+fieldTextWidth > interiorWidth {
		adjustedFieldCol = interiorWidth - fieldTextWidth
	}

	// Ensure field doesn't start before left edge
	if adjustedFieldCol < 0 {
		adjustedFieldCol = 0
	}

	// Store adjusted field position (convert back to screen coords)
	m.fieldRow = fieldRow
	m.fieldCol = adjustedFieldCol + 1 // +1 to convert back to screen coords

	// Create dropdown - pass adjusted field position
	// Set margins to exclude menu bar (top) and status bar (bottom)
	m.dropdown = teadrpdwn.NewDropdownModel(items, &teadrpdwn.DropdownModelArgs{
		FieldRow:     m.fieldRow,
		FieldCol:     m.fieldCol,
		ScreenWidth:  m.screenWidth,
		ScreenHeight: m.screenHeight,
		TopMargin:    1, // Menu bar at row 0
		BottomMargin: 1, // Status bar at row screenHeight-1
	})

	if m.selectedStartOpen {
		m.dropdown, cmd = m.dropdown.Open()
		_ = cmd // Ignore cmd since we're not in Update()
	}

	return m
}

func (m exampleModel) View() tea.View {
	var content string
	if m.mode == modeConfiguring {
		content = m.viewConfiguring()
	} else {
		content = m.viewDemonstrating()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m exampleModel) viewConfiguring() string {
	var lines []string

	// Row 0: Title centered
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Width(m.screenWidth).
		Align(lipgloss.Center).
		Render("Configure Dropdown Demo")
	lines = append(lines, title)

	// Blank line
	lines = append(lines, "")

	// Render each field centered in 40-char box
	horizontalStr := ""
	if m.configComplete[0] {
		horizontalStr = m.horizontalOptions[m.selectedHorizontal]
	}
	lines = append(lines, m.renderFieldCentered(0, "Horizontal Position", horizontalStr, m.renderSelector(m.horizontalOptions, m.horizontalIdx))...)

	verticalStr := ""
	if m.configComplete[1] {
		verticalStr = m.verticalOptions[m.selectedVertical]
	}
	lines = append(lines, m.renderFieldCentered(1, "Vertical Position", verticalStr, m.renderSelector(m.verticalOptions, m.verticalIdx))...)

	numItemsStr := ""
	if m.configComplete[2] {
		numItemsStr = fmt.Sprintf("%d items", m.selectedNumItems)
	}
	lines = append(lines, m.renderFieldCentered(2, "Number of Items", numItemsStr, m.renderSelector(m.numItemsOptions, m.numItemsIdx))...)

	widthStr := ""
	if m.configComplete[3] {
		widthStr = m.getWidthDescription()
	}
	lines = append(lines, m.renderFieldCentered(3, "Width of Longest Item", widthStr, m.renderSelector(m.widthOptions, m.widthIdx))...)

	startOpenStr := ""
	if m.configComplete[4] {
		if m.selectedStartOpen {
			startOpenStr = "Yes"
		} else {
			startOpenStr = "No"
		}
	}
	lines = append(lines, m.renderFieldCentered(4, "Start Open", startOpenStr, m.renderSelector(m.startOpenOptions, m.startOpenIdx))...)

	// Bottom: Help text
	for len(lines) < m.screenHeight-1 {
		lines = append(lines, "")
	}
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Background(lipgloss.Color("235")).
		Width(m.screenWidth).
		Align(lipgloss.Center).
		Render("↑↓ navigate | Enter confirm | Esc back | q quit")
	lines = append(lines, help)

	return strings.Join(lines, "\n")
}

func (m exampleModel) renderSelector(options []string, selectedIdx int) string {
	var b strings.Builder
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("170"))

	for i, option := range options {
		if i == selectedIdx {
			b.WriteString(selectedStyle.Render("> " + option))
		} else {
			b.WriteString("  " + option)
		}
		if i < len(options)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m exampleModel) renderFieldCentered(fieldNum int, label, value, selectorView string) []string {
	var lines []string
	var symbol string
	const boxWidth = 50 // Increased for longer labels

	if m.configComplete[fieldNum] {
		symbol = "✓"
	} else if m.currentField == fieldNum {
		symbol = "→"
	} else {
		symbol = "○"
	}

	// Calculate centering offset
	leftPad := (m.screenWidth - boxWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	if m.configComplete[fieldNum] {
		line := fmt.Sprintf("%s %s: %s", symbol, label, value)
		lines = append(lines, strings.Repeat(" ", leftPad)+line)
	} else if m.currentField == fieldNum {
		// Header
		header := fmt.Sprintf("%s %s", symbol, label)
		lines = append(lines, strings.Repeat(" ", leftPad)+header)
		// Selector lines
		selectorLines := strings.Split(selectorView, "\n")
		for _, sline := range selectorLines {
			lines = append(lines, strings.Repeat(" ", leftPad)+sline)
		}
	} else {
		line := fmt.Sprintf("%s %s: (not selected)", symbol, label)
		lines = append(lines, strings.Repeat(" ", leftPad)+line)
	}

	// No blank line after field (removed to save vertical space)

	return lines
}

func (m exampleModel) viewDemonstrating() string {
	// Build the screen content as a 2D structure that we'll render at the end

	// Menu bar at top
	menu := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Width(m.screenWidth).
		Align(lipgloss.Center).
		Render("Space/Enter/Down/Right: open | Esc/r: reconfigure | q: quit")

	// Config status at bottom
	startOpenStr := "No"
	if m.selectedStartOpen {
		startOpenStr = "Yes"
	}
	configStatus := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Background(lipgloss.Color("235")).
		Width(m.screenWidth).
		Align(lipgloss.Center).
		Render(fmt.Sprintf(
			"H=%s V=%s | Items=%d | Width=%s | StartOpen=%s",
			m.horizontalOptions[m.selectedHorizontal],
			m.verticalOptions[m.selectedVertical],
			m.selectedNumItems,
			m.getWidthDescription(),
			startOpenStr,
		))

	// Build the bordered screen area with field
	screenArea := m.buildScreenArea()

	// Overlay dropdown if open
	if m.dropdown.IsOpen {
		dropdownView := m.dropdown.View().Content
		// Adjust dropdown row to be relative to screenArea (which starts at screen row 1)
		// Menu bar is at screen row 0, screenArea starts at screen row 1
		relativeRow := m.dropdown.Row - 1
		screenArea = teadrpdwn.OverlayDropdown(screenArea, dropdownView, relativeRow, m.dropdown.Col)
	}

	// Combine vertically
	return lipgloss.JoinVertical(lipgloss.Left,
		menu,
		screenArea,
		configStatus,
	)
}

func (m exampleModel) buildScreenArea() string {
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(m.screenWidth).
		Height(m.screenHeight - 2) // Account for menu and status bars

	// Build the content for inside the border
	content := m.buildScreenContent()

	return borderStyle.Render(content)
}

func (m exampleModel) buildScreenContent() string {
	interiorWidth := m.screenWidth - 2
	interiorHeight := m.screenHeight - 3

	// Use stored field position (independent of dropdown position)
	fieldRow := m.fieldRow - 2 // Adjust for menu and border
	fieldCol := m.fieldCol - 1 // Adjust for left border

	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("170")).
		Bold(true)

	var lines []string
	for row := 0; row < interiorHeight; row++ {
		if row == fieldRow {
			// Render the field line
			triangle := "▶"
			if m.dropdown.IsOpen {
				if m.dropdown.DisplayAbove {
					triangle = "▲"
				} else {
					triangle = "▼"
				}
			}
			fieldText := triangle + " " + m.selectedValue
			fieldTextWidth := ansi.StringWidth(fieldText)

			// Truncate if field text is wider than available space
			if fieldTextWidth > interiorWidth {
				fieldText = ansi.Truncate(fieldText, interiorWidth-1, "…")
			}

			// Build plain text line with proper spacing (fieldCol already adjusted in setupDemo)
			plainLine := strings.Repeat(" ", fieldCol) + fieldText
			// Pad to full width using ANSI-aware width
			if ansi.StringWidth(plainLine) < interiorWidth {
				plainLine = plainLine + strings.Repeat(" ", interiorWidth-ansi.StringWidth(plainLine))
			}
			// Apply styling
			lines = append(lines, fieldStyle.Render(plainLine))
		} else {
			// Empty line
			lines = append(lines, strings.Repeat(" ", interiorWidth))
		}
	}

	return strings.Join(lines, "\n")
}

func generateItems(widthType string, count, screenWidth int) []teadrpdwn.Option {
	var allItems []string

	switch widthType {
	case "short":
		// Simple short items
		allItems = []string{
			"Red", "Green", "Blue", "Yellow", "Orange",
			"Purple", "Pink", "Brown", "Gray", "Black",
			"White", "Cyan", "Magenta", "Lime", "Navy",
			"Teal", "Olive", "Maroon", "Aqua", "Silver",
			"Gold", "Indigo", "Violet", "Coral", "Salmon",
		}
	case "medium":
		// Medium-length items (file paths, descriptions)
		allItems = []string{
			"src/components/Button.tsx",
			"src/components/Dropdown.tsx",
			"config/database.yaml",
			"tests/integration/api_test.go",
			"docs/architecture/overview.md",
			"internal/handlers/user.go",
			"pkg/utils/validation.go",
			"cmd/server/main.go",
			"assets/images/logo.png",
			"scripts/deploy.sh",
			"migrations/001_initial.sql",
			"locales/en-US/messages.json",
			"templates/email/welcome.html",
			"public/static/css/main.css",
			"api/v1/endpoints/auth.go",
			"models/repositories/user_repo.go",
			"middleware/logging/logger.go",
			"services/payment/stripe.go",
			"workers/background/email_worker.go",
			"infrastructure/terraform/main.tf",
			"monitoring/prometheus/metrics.go",
			"cache/redis/connection.go",
			"queue/rabbitmq/consumer.go",
			"storage/s3/uploader.go",
			"auth/jwt/validator.go",
		}
	case "long":
		// Long items (deeply nested paths - will demonstrate truncation)
		allItems = []string{
			"/var/www/application/storage/framework/sessions/cache/data/production/user_sessions_2024_archive",
			"/usr/local/share/applications/development/tools/jetbrains/goland/plugins/configuration/settings",
			"/opt/homebrew/Cellar/go/1.21.5/libexec/src/runtime/internal/atomic/atomic_amd64.s",
			"/Users/developer/Projects/company-monorepo/packages/backend/services/api/handlers/authentication",
			"/System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework",
			"/var/storage/cloud-backups/bucket-name/uploads/users/profile-images/thumbnails/original/metadata",
			"/opt/applications/microservices/user-service/internal/domain/repositories/implementations/",
			"/usr/share/documentation/guides/getting-started/installation/prerequisites/system-requirements",
			"/home/developer/workspace/frontend-application/node_modules/@company/shared-components/dist/",
			"/var/cache/npm/registry/packages/@organization/package-name/versions/package-name-1.2.3.tgz",
			"/Library/Application Support/Company/Product/Configuration/Environments/Production/Settings",
			"/var/data/marketplace/api/v3/products/categories/electronics/computers/laptops/specifications",
			"/usr/local/share/applications/development-tools/ide/plugins/extensions/community/",
			"/opt/analytics/dashboard/reports/2024/quarterly/revenue/breakdown/details/summaries",
			"/System/Library/PrivateFrameworks/CoreServices.framework/Versions/A/Resources/",
			"/var/ecommerce/store/checkout/payment/methods/credit-card/validation/secure/tokens",
			"/var/log/applications/production/web-servers/nginx/access-logs/2024/january/daily/",
			"/var/mail/archives/inbox/messages/2024/january/processed/attachments/documents",
			"/Applications/Development Tools/IDEs/IntelliJ IDEA.app/Contents/Resources/",
			"/mnt/network/shared/drive/folders/projects/2024/client-work/deliverables/final/approved",
			"/mnt/storage/backups/databases/production/daily/compressed/encrypted/archived/",
			"/var/forums/questions/12345678/answers/accepted/how-to-implement-complex-feature-with-multiple",
			"/private/var/folders/xy/abcdefgh12345678/T/TemporaryItems/com.company.product/",
			"/var/media/video/content/library/collections/playlist-items/metadata/thumbnails/high-resolution",
			"/Users/Shared/Development/Repositories/opensource-projects/contributions/pull-requests/",
		}
	default:
		allItems = []string{"Item 1", "Item 2", "Item 3"}
	}

	// Return the requested number of items
	if count > len(allItems) {
		count = len(allItems)
	}
	return teadrpdwn.ToOptions(allItems[:count])
}

func calculateFieldPosition(horizontal, vertical, width, height int) (row, col int) {
	// Border coordinates: top-left=(1,0), bottom-right=(height-2, width-1)
	// Interior starts at (2, 1)
	borderTop := 2
	borderLeft := 1
	borderRight := width - 2
	borderBottom := height - 3

	borderHeight := borderBottom - borderTop
	borderWidth := borderRight - borderLeft

	// Calculate row based on vertical position (0=Top, 1=Middle, 2=Bottom)
	switch vertical {
	case 0: // Top
		row = borderTop
	case 1: // Middle
		row = borderTop + borderHeight/2
	case 2: // Bottom
		row = borderBottom - 2
	default:
		row = borderTop
	}

	// Calculate col based on horizontal position (0=Left, 1=Middle, 2=Right)
	switch horizontal {
	case 0: // Left
		col = borderLeft + 3
	case 1: // Middle
		col = borderLeft + borderWidth/2
	case 2: // Right
		col = borderRight - 20
	default:
		col = borderLeft + 3
	}

	return row, col
}
