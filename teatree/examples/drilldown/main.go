package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

// Dataset names and number key bindings
var datasetNames = []string{
	"File System",
	"Product Taxonomy",
	"Geography",
	"Org Chart",
	"Menu Hierarchy",
}

type model struct {
	drillDown     teatree.DrillDownModel[string]
	datasets      []*teatree.Node[string]
	activeDataset int
	width         int
	height        int
	statusMsg     string
}

func initialModel() model {
	datasets := buildAllDatasets()

	m := model{
		datasets:      datasets,
		activeDataset: 0,
	}
	m.drillDown = newDrillDown(datasets[0])
	return m
}

func newDrillDown(root *teatree.Node[string]) teatree.DrillDownModel[string] {
	dd := teatree.NewDrillDownModel(root, teatree.DrillDownArgs[string]{
		SelectorFunc: firstChildSelector,
		Prompt:       "Drill Down Explorer",
	})
	dd, _ = dd.Initialize()
	dd.InsertLine = true
	return dd
}

func firstChildSelector(_ *teatree.Node[string], children []*teatree.Node[string]) (*teatree.Node[string], error) {
	if len(children) == 0 {
		return nil, nil
	}
	return children[0], nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.drillDown.Width = msg.Width - 4
		m.drillDown.Height = msg.Height - 6

	case tea.KeyPressMsg:
		// Quit on q/ctrl+c
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Number keys 1-5 switch datasets
		if msg.Text >= "1" && msg.Text <= "5" {
			idx := int(msg.Text[0]-'0') - 1
			if idx != m.activeDataset {
				m.activeDataset = idx
				m.drillDown = newDrillDown(m.datasets[idx])
				m.drillDown.Width = m.width - 4
				m.drillDown.Height = m.height - 6
				m.statusMsg = fmt.Sprintf("Switched to: %s", datasetNames[idx])
				return m, nil
			}
		}

	case teatree.DrillDownSelectMsg[string]:
		m.statusMsg = fmt.Sprintf("Selected: %s — %s", msg.Node.Name(), *msg.Node.Data())
		return m, nil

	case teatree.DrillDownChangeMsg[string]:
		m.statusMsg = fmt.Sprintf("Changed level %d to: %s", msg.Level, msg.Node.Name())
		return m, nil

	case teatree.DrillDownFocusMsg[string]:
		m.statusMsg = fmt.Sprintf("Focused: %s (level %d)", msg.Node.Name(), msg.Level)
		return m, nil
	}

	// Delegate to drill-down model
	var ddModel tea.Model
	ddModel, cmd = m.drillDown.Update(msg)
	m.drillDown = ddModel.(teatree.DrillDownModel[string])

	return m, cmd
}

func (m model) View() tea.View {
	var view string

	// Header with dataset tabs
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("Drill Down Demo")
	header += "  "
	for i, name := range datasetNames {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		if i == m.activeDataset {
			style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("230")).Background(lipgloss.Color("62"))
		}
		header += style.Render(fmt.Sprintf(" %d:%s ", i+1, name))
		if i < len(datasetNames)-1 {
			header += " "
		}
	}

	// Drill-down view
	ddView := m.drillDown.View().Content

	// Status line
	status := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(m.statusMsg)

	// Footer
	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("↑↓ navigate • space/→ alternatives • enter select • 1-5 dataset • q quit")

	view = header + "\n\n" + ddView + "\n\n" + status + "\n" + footer

	return tea.NewView(view)
}

// =============================================================================
// Datasets
// =============================================================================

func buildAllDatasets() []*teatree.Node[string] {
	return []*teatree.Node[string]{
		buildFileSystemTree(),
		buildProductTree(),
		buildGeographyTree(),
		buildOrgChartTree(),
		buildMenuTree(),
	}
}

func buildFileSystemTree() *teatree.Node[string] {
	root := teatree.NewNode("fs", "~/Projects", "Home projects directory")

	goProjects := teatree.NewNode("go", "go-pkgs", "Go packages")
	webProjects := teatree.NewNode("web", "web-apps", "Web applications")
	docsDir := teatree.NewNode("docs", "docs", "Documentation")

	goDt := teatree.NewNode("go-dt", "go-dt", "Domain types package")
	goTest := teatree.NewNode("go-test", "go-testutil", "Test utilities")
	goTealeaves := teatree.NewNode("go-tl", "go-tealeaves", "TUI components")

	react := teatree.NewNode("react", "react-dashboard", "React dashboard app")
	vue := teatree.NewNode("vue", "vue-storefront", "Vue storefront")

	readme := teatree.NewNode("readme", "README.md", "Project readme")
	changelog := teatree.NewNode("changelog", "CHANGELOG.md", "Change log")

	root.AddChild(goProjects)
	root.AddChild(webProjects)
	root.AddChild(docsDir)

	goProjects.AddChild(goDt)
	goProjects.AddChild(goTest)
	goProjects.AddChild(goTealeaves)

	webProjects.AddChild(react)
	webProjects.AddChild(vue)

	docsDir.AddChild(readme)
	docsDir.AddChild(changelog)

	return root
}

func buildProductTree() *teatree.Node[string] {
	root := teatree.NewNode("store", "Store", "Online retail store")

	electronics := teatree.NewNode("elec", "Electronics", "Electronic devices")
	clothing := teatree.NewNode("cloth", "Clothing", "Apparel and accessories")
	books := teatree.NewNode("books", "Books", "Books and media")

	phones := teatree.NewNode("phones", "Phones", "Mobile phones")
	laptops := teatree.NewNode("laptops", "Laptops", "Laptop computers")

	iphone := teatree.NewNode("iphone", "iPhone 16", "Apple iPhone 16")
	pixel := teatree.NewNode("pixel", "Pixel 9", "Google Pixel 9")
	galaxy := teatree.NewNode("galaxy", "Galaxy S25", "Samsung Galaxy S25")

	macbook := teatree.NewNode("macbook", "MacBook Pro", "Apple MacBook Pro")
	thinkpad := teatree.NewNode("thinkpad", "ThinkPad X1", "Lenovo ThinkPad X1")

	mens := teatree.NewNode("mens", "Men's", "Men's clothing")
	womens := teatree.NewNode("womens", "Women's", "Women's clothing")

	fiction := teatree.NewNode("fiction", "Fiction", "Fiction books")
	tech := teatree.NewNode("tech", "Technology", "Technology books")

	root.AddChild(electronics)
	root.AddChild(clothing)
	root.AddChild(books)

	electronics.AddChild(phones)
	electronics.AddChild(laptops)

	phones.AddChild(iphone)
	phones.AddChild(pixel)
	phones.AddChild(galaxy)

	laptops.AddChild(macbook)
	laptops.AddChild(thinkpad)

	clothing.AddChild(mens)
	clothing.AddChild(womens)

	books.AddChild(fiction)
	books.AddChild(tech)

	return root
}

func buildGeographyTree() *teatree.Node[string] {
	root := teatree.NewNode("world", "World", "Planet Earth")

	usa := teatree.NewNode("usa", "United States", "USA")
	canada := teatree.NewNode("canada", "Canada", "Canada")
	uk := teatree.NewNode("uk", "United Kingdom", "UK")

	california := teatree.NewNode("ca", "California", "The Golden State")
	texas := teatree.NewNode("tx", "Texas", "The Lone Star State")
	newYork := teatree.NewNode("ny", "New York", "The Empire State")

	sf := teatree.NewNode("sf", "San Francisco", "City by the Bay")
	la := teatree.NewNode("la", "Los Angeles", "City of Angels")
	sd := teatree.NewNode("sd", "San Diego", "America's Finest City")

	houston := teatree.NewNode("hou", "Houston", "Space City")
	austin := teatree.NewNode("aus", "Austin", "Live Music Capital")

	nyc := teatree.NewNode("nyc", "New York City", "The Big Apple")
	buffalo := teatree.NewNode("buf", "Buffalo", "The Nickel City")

	ontario := teatree.NewNode("on", "Ontario", "Province of Ontario")
	bc := teatree.NewNode("bc", "British Columbia", "Beautiful BC")

	toronto := teatree.NewNode("tor", "Toronto", "The Six")
	vancouver := teatree.NewNode("van", "Vancouver", "Raincouver")

	england := teatree.NewNode("eng", "England", "England")
	london := teatree.NewNode("lon", "London", "The Big Smoke")

	root.AddChild(usa)
	root.AddChild(canada)
	root.AddChild(uk)

	usa.AddChild(california)
	usa.AddChild(texas)
	usa.AddChild(newYork)

	california.AddChild(sf)
	california.AddChild(la)
	california.AddChild(sd)

	texas.AddChild(houston)
	texas.AddChild(austin)

	newYork.AddChild(nyc)
	newYork.AddChild(buffalo)

	canada.AddChild(ontario)
	canada.AddChild(bc)

	ontario.AddChild(toronto)
	bc.AddChild(vancouver)

	uk.AddChild(england)
	england.AddChild(london)

	return root
}

func buildOrgChartTree() *teatree.Node[string] {
	root := teatree.NewNode("co", "Acme Corp", "Global technology company")

	eng := teatree.NewNode("eng", "Engineering", "Software engineering division")
	sales := teatree.NewNode("sales", "Sales", "Sales and business development")
	ops := teatree.NewNode("ops", "Operations", "IT and infrastructure")

	frontend := teatree.NewNode("fe", "Frontend", "UI/UX engineering team")
	backend := teatree.NewNode("be", "Backend", "Server-side engineering team")
	infra := teatree.NewNode("infra", "Infrastructure", "Platform and DevOps team")

	alice := teatree.NewNode("alice", "Alice Chen", "Senior Frontend Engineer")
	bob := teatree.NewNode("bob", "Bob Smith", "Frontend Engineer")

	charlie := teatree.NewNode("charlie", "Charlie Park", "Senior Backend Engineer")
	diana := teatree.NewNode("diana", "Diana Lee", "Backend Engineer")

	eve := teatree.NewNode("eve", "Eve Johnson", "SRE Lead")

	frank := teatree.NewNode("frank", "Frank Williams", "Enterprise Sales")
	grace := teatree.NewNode("grace", "Grace Kim", "SMB Sales")

	root.AddChild(eng)
	root.AddChild(sales)
	root.AddChild(ops)

	eng.AddChild(frontend)
	eng.AddChild(backend)
	eng.AddChild(infra)

	frontend.AddChild(alice)
	frontend.AddChild(bob)

	backend.AddChild(charlie)
	backend.AddChild(diana)

	infra.AddChild(eve)

	sales.AddChild(frank)
	sales.AddChild(grace)

	return root
}

func buildMenuTree() *teatree.Node[string] {
	root := teatree.NewNode("main", "Main Menu", "Application main menu")

	file := teatree.NewNode("file", "File", "File operations")
	edit := teatree.NewNode("edit", "Edit", "Edit operations")
	view := teatree.NewNode("view", "View", "View settings")
	help := teatree.NewNode("help", "Help", "Help and documentation")

	newFile := teatree.NewNode("new", "New", "Create new document")
	open := teatree.NewNode("open", "Open", "Open existing document")
	save := teatree.NewNode("save", "Save", "Save current document")
	export := teatree.NewNode("export", "Export", "Export to other formats")

	pdf := teatree.NewNode("pdf", "PDF", "Export as PDF")
	csv := teatree.NewNode("csv", "CSV", "Export as CSV")
	html := teatree.NewNode("html", "HTML", "Export as HTML")

	undo := teatree.NewNode("undo", "Undo", "Undo last action")
	redo := teatree.NewNode("redo", "Redo", "Redo last undone action")
	find := teatree.NewNode("find", "Find & Replace", "Search and replace text")

	zoom := teatree.NewNode("zoom", "Zoom", "Zoom controls")
	theme := teatree.NewNode("theme", "Theme", "Appearance settings")

	zoomIn := teatree.NewNode("zoomin", "Zoom In", "Increase zoom level")
	zoomOut := teatree.NewNode("zoomout", "Zoom Out", "Decrease zoom level")

	about := teatree.NewNode("about", "About", "About this application")
	docs := teatree.NewNode("docs", "Documentation", "Open documentation")

	root.AddChild(file)
	root.AddChild(edit)
	root.AddChild(view)
	root.AddChild(help)

	file.AddChild(newFile)
	file.AddChild(open)
	file.AddChild(save)
	file.AddChild(export)

	export.AddChild(pdf)
	export.AddChild(csv)
	export.AddChild(html)

	edit.AddChild(undo)
	edit.AddChild(redo)
	edit.AddChild(find)

	view.AddChild(zoom)
	view.AddChild(theme)

	zoom.AddChild(zoomIn)
	zoom.AddChild(zoomOut)

	help.AddChild(about)
	help.AddChild(docs)

	return root
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
