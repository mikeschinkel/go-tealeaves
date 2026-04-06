package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	colID      = "id"
	colFirst   = "first"
	colLast    = "last"
	colDept    = "dept"
	colTitle   = "title"
	colEmail   = "email"
	colPhone   = "phone"
	colOffice  = "office"
	colFloor   = "floor"
	colStart   = "start"
	colSalary  = "salary"
	colManager = "manager"
	colStatus  = "status"
	colCity    = "city"
	colCountry = "country"
)

type employee struct {
	id, first, last, dept, title   string
	email, phone, office, floor    string
	start, salary, manager, status string
	city, country                  string
}

var employees = []employee{
	{"E001", "Alice", "Chen", "Engineering", "Senior Engineer", "achen@example.com", "+1-555-0101", "B3-401", "4", "2019-03-15", "$145,000", "Bob Torres", "Active", "San Francisco", "USA"},
	{"E002", "Bob", "Torres", "Engineering", "Engineering Manager", "btorres@example.com", "+1-555-0102", "B3-400", "4", "2017-06-01", "$175,000", "VP Eng", "Active", "San Francisco", "USA"},
	{"E003", "Christopher", "Vandenberg", "Customer Success", "CS Lead", "cvandenberg@example.com", "+44-20-7946-0958", "L2-210", "2", "2020-01-10", "$120,000", "Diana Osei", "Active", "London", "United Kingdom"},
	{"E004", "Diana", "Osei", "Customer Success", "VP Customer Success", "dosei@example.com", "+44-20-7946-0959", "L2-200", "2", "2018-09-22", "$165,000", "CEO", "Active", "London", "United Kingdom"},
	{"E005", "Li", "Wei", "Sales", "Account Executive", "lwei@example.com", "+1-555-0105", "B1-105", "1", "2021-07-14", "$95,000", "Frank Müller", "Active", "New York", "USA"},
	{"E006", "Frank", "Müller", "Sales", "Sales Director", "fmuller@example.com", "+49-30-1234-5678", "B1-100", "1", "2016-11-03", "$155,000", "CEO", "Active", "Berlin", "Germany"},
	{"E007", "Priya", "Sharma", "Engineering", "Staff Engineer", "psharma@example.com", "+91-22-2345-6789", "R4-310", "3", "2018-02-28", "$160,000", "Bob Torres", "Active", "Mumbai", "India"},
	{"E008", "James", "O'Brien", "Finance", "Financial Analyst", "jobrien@example.com", "+1-555-0108", "B2-205", "2", "2022-04-18", "$88,000", "Sara Kim", "Active", "Chicago", "USA"},
	{"E009", "Sara", "Kim", "Finance", "CFO", "skim@example.com", "+1-555-0109", "B2-200", "2", "2015-08-07", "$195,000", "CEO", "Active", "Chicago", "USA"},
	{"E010", "Yuki", "Tanaka", "Engineering", "DevOps Engineer", "ytanaka@example.com", "+81-3-1234-5678", "T1-501", "5", "2020-11-30", "$130,000", "Bob Torres", "Active", "Tokyo", "Japan"},
	{"E011", "Maria", "Garcia", "Marketing", "Content Strategist", "mgarcia@example.com", "+34-91-123-4567", "M3-302", "3", "2021-03-22", "$92,000", "Tom Novak", "On Leave", "Madrid", "Spain"},
	{"E012", "Tom", "Novak", "Marketing", "VP Marketing", "tnovak@example.com", "+1-555-0112", "B1-200", "2", "2017-01-15", "$170,000", "CEO", "Active", "New York", "USA"},
	{"E013", "Aisha", "Patel", "Engineering", "QA Lead", "apatel@example.com", "+1-555-0113", "B3-405", "4", "2019-09-01", "$125,000", "Bob Torres", "Active", "San Francisco", "USA"},
	{"E014", "Erik", "Johansson", "Support", "Support Manager", "ejohansson@example.com", "+46-8-123-4567", "S2-101", "1", "2020-06-15", "$105,000", "Diana Osei", "Active", "Stockholm", "Sweden"},
	{"E015", "Fatima", "Al-Rashid", "Engineering", "Backend Engineer", "falrashid@example.com", "+971-4-123-4567", "D1-401", "4", "2022-01-10", "$135,000", "Bob Torres", "Active", "Dubai", "UAE"},
	{"E016", "Carlos", "Silva", "Sales", "Sales Rep", "csilva@example.com", "+55-11-1234-5678", "SP-105", "1", "2023-02-01", "$78,000", "Frank Müller", "Active", "São Paulo", "Brazil"},
	{"E017", "Hannah", "Williams", "HR", "HR Director", "hwilliams@example.com", "+1-555-0117", "B2-300", "3", "2018-04-10", "$140,000", "CEO", "Active", "San Francisco", "USA"},
	{"E018", "Raj", "Krishnamurthy", "Engineering", "Frontend Engineer", "rkrishnamurthy@example.com", "+91-80-2345-6789", "R4-315", "3", "2021-08-20", "$115,000", "Bob Torres", "Active", "Bangalore", "India"},
	{"E019", "Sophie", "Dubois", "Legal", "General Counsel", "sdubois@example.com", "+33-1-2345-6789", "P1-200", "2", "2019-05-15", "$185,000", "CEO", "Active", "Paris", "France"},
	{"E020", "Oleksandr", "Kovalenko", "Engineering", "SRE", "okovalenko@example.com", "+48-22-123-4567", "W1-401", "4", "2022-09-01", "$128,000", "Bob Torres", "Remote", "Warsaw", "Poland"},
}

type colDef struct {
	key   string
	title string
	field func(employee) string
}

var colDefs = []colDef{
	{colID, "ID", func(e employee) string { return e.id }},
	{colFirst, "First", func(e employee) string { return e.first }},
	{colLast, "Last", func(e employee) string { return e.last }},
	{colDept, "Dept", func(e employee) string { return e.dept }},
	{colTitle, "Title", func(e employee) string { return e.title }},
	{colEmail, "Email", func(e employee) string { return e.email }},
	{colPhone, "Phone", func(e employee) string { return e.phone }},
	{colOffice, "Office", func(e employee) string { return e.office }},
	{colManager, "Manager", func(e employee) string { return e.manager }},
	{colFloor, "Floor", func(e employee) string { return e.floor }},
	{colStart, "Start", func(e employee) string { return e.start }},
	{colSalary, "Salary", func(e employee) string { return e.salary }},
	{colStatus, "Status", func(e employee) string { return e.status }},
	{colCity, "City", func(e employee) string { return e.city }},
	{colCountry, "Country", func(e employee) string { return e.country }},
}

func maxLen(header string, field func(employee) string) int {
	n := len(header)
	for _, e := range employees {
		if l := len(field(e)); l > n {
			n = l
		}
	}
	return n
}

type model struct {
	table       teagrid.GridModel
	colDefs     []colDef
	width       int
	height      int
	freezeCount int
	colCursor   bool
	wrapping    bool
	overflow    bool
}

func newModel() model {
	columns := make([]teagrid.Column, len(colDefs))
	for i, cd := range colDefs {
		switch cd.key {
		case colEmail:
			columns[i] = teagrid.NewFlexColumn(cd.key, cd.title, 2)
		case colCountry:
			columns[i] = teagrid.NewFlexColumn(cd.key, cd.title, 1)
		case colID:
			columns[i] = teagrid.NewColumn(cd.key, cd.title, maxLen(cd.title, cd.field)).
				WithAlignment(lipgloss.Center)
		case colSalary:
			columns[i] = teagrid.NewColumn(cd.key, cd.title, maxLen(cd.title, cd.field)).
				WithAlignment(lipgloss.Right)
		default:
			columns[i] = teagrid.NewColumn(cd.key, cd.title, maxLen(cd.title, cd.field))
		}
	}

	rows := make([]teagrid.Row, len(employees))
	for i, e := range employees {
		rows[i] = teagrid.NewRow(teagrid.RowData{
			colID:      e.id,
			colFirst:   e.first,
			colLast:    e.last,
			colDept:    e.dept,
			colTitle:   e.title,
			colEmail:   e.email,
			colPhone:   e.phone,
			colOffice:  e.office,
			colFloor:   e.floor,
			colStart:   e.start,
			colSalary:  e.salary,
			colManager: e.manager,
			colStatus:  e.status,
			colCity:    e.city,
			colCountry: e.country,
		})
	}

	freezeCount := 2
	const colCursor = true
	const overflow = true

	return model{
		table: teagrid.
			NewGridModel(columns).
			WithRows(rows).
			WithHorizontalFreezeColumnCount(freezeCount).
			WithColCursorMode(colCursor).
			WithOverflowIndicator(overflow).
			WithFillWidth(true).
			WithFocused(true),
		colDefs:     colDefs,
		freezeCount: freezeCount,
		colCursor:   colCursor,
		overflow:    overflow,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Subtract 4 lines: title + debug status + blank line + help line
		m.table = m.table.WithSize(msg.Width, msg.Height-4)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)

		case "f":
			m.freezeCount = (m.freezeCount + 1) % 4
			m.table = m.table.WithHorizontalFreezeColumnCount(m.freezeCount)

		case "c":
			m.colCursor = !m.colCursor
			m.table = m.table.WithColCursorMode(m.colCursor)

		case "w":
			m.wrapping = !m.wrapping
			m.table = m.table.
				WithRowCursorWrapping(m.wrapping).
				WithColCursorWrapping(m.wrapping)

		case "o":
			m.overflow = !m.overflow
			m.table = m.table.WithOverflowIndicator(m.overflow)

		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) debugStatusLine() string {
	cursorCol := m.table.ColCursorColumnIndex()
	colName := ""
	if cursorCol < len(m.colDefs) {
		colName = m.colDefs[cursorCol].title
	}

	colMode := "off"
	if m.colCursor {
		colMode = "on"
	}

	wrapMode := "off"
	if m.wrapping {
		wrapMode = "on"
	}

	overflowMode := "off"
	if m.overflow {
		overflowMode = "on"
	}

	return fmt.Sprintf("[%dx%d] Scroll:%d Cursor:%d=%q Freeze:%d ColCursor:%s Wrap:%s Overflow:%s Page:%d/%d Rows:%d Total:%dw Natural:%dw",
		m.width, m.height,
		m.table.HorizontalScrollColumnOffset(),
		cursorCol, colName,
		m.freezeCount,
		colMode,
		wrapMode,
		overflowMode,
		m.table.CurrentPage(), m.table.MaxPages(),
		m.table.TotalRows(),
		m.table.TotalWidth(), m.table.NaturalWidth(),
	)
}

func (m model) View() tea.View {
	var body strings.Builder

	body.WriteString("Employee Directory — Horizontal Panning Demo\n")
	body.WriteString(m.debugStatusLine())
	body.WriteString("\n\n")
	body.WriteString(m.table.View().Content)

	helpText := "f: freeze | c: col cursor | w: wrap | o: overflow | q: quit"
	helpWidth := len(helpText)
	if pad := (m.width - helpWidth) / 2; pad > 0 {
		body.WriteString("\n" + strings.Repeat(" ", pad) + helpText)
	} else {
		body.WriteString("\n" + helpText)
	}

	v := tea.NewView(body.String())
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
