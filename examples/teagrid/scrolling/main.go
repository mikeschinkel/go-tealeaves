package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
)

const (
	colName   = "name"
	colIP     = "ip"
	colOS     = "os"
	colCPU    = "cpu"
	colRAM    = "ram"
	colDisk   = "disk"
	colStatus = "status"
	colRegion = "region"
	colUptime = "uptime"
	colOwner  = "owner"
)

type server struct {
	name   string
	ip     string
	os     string
	cpu    string
	ram    string
	disk   string
	status string
	region string
	uptime string
	owner  string
}

var servers = []server{
	{"web-prod-01", "10.0.1.12", "Ubuntu 22.04", "4 vCPU", "16 GB", "200 GB SSD", "running", "us-east-1", "47d 12h", "Platform Team"},
	{"web-prod-02", "10.0.1.13", "Ubuntu 22.04", "4 vCPU", "16 GB", "200 GB SSD", "running", "us-east-1", "47d 12h", "Platform Team"},
	{"api-prod-01", "10.0.2.20", "Debian 12", "8 vCPU", "32 GB", "500 GB SSD", "running", "us-west-2", "12d 3h", "Backend"},
	{"api-prod-02", "10.0.2.21", "Debian 12", "8 vCPU", "32 GB", "500 GB SSD", "degraded", "us-west-2", "12d 3h", "Backend"},
	{"db-primary", "10.0.3.10", "Rocky Linux 9", "16 vCPU", "128 GB", "2 TB NVMe", "running", "us-east-1", "90d 5h", "DBA Team"},
	{"db-replica-01", "10.0.3.11", "Rocky Linux 9", "16 vCPU", "128 GB", "2 TB NVMe", "running", "us-east-1", "90d 5h", "DBA Team"},
	{"db-replica-02", "10.0.3.12", "Rocky Linux 9", "8 vCPU", "64 GB", "1 TB NVMe", "running", "eu-west-1", "45d 8h", "DBA Team"},
	{"cache-01", "10.0.4.5", "Alpine 3.19", "2 vCPU", "8 GB", "50 GB SSD", "running", "us-east-1", "30d 1h", "Platform Team"},
	{"cache-02", "10.0.4.6", "Alpine 3.19", "2 vCPU", "8 GB", "50 GB SSD", "stopped", "us-west-2", "0d 0h", "Platform Team"},
	{"worker-01", "10.0.5.30", "Ubuntu 24.04", "32 vCPU", "256 GB", "4 TB NVMe", "running", "us-east-1", "5d 22h", "ML Team"},
	{"worker-02", "10.0.5.31", "Ubuntu 24.04", "32 vCPU", "256 GB", "4 TB NVMe", "running", "us-east-1", "5d 22h", "ML Team"},
	{"monitor", "10.0.6.2", "Debian 12", "2 vCPU", "4 GB", "100 GB SSD", "running", "us-east-1", "120d 0h", "SRE"},
	{"bastion", "10.0.0.5", "Ubuntu 22.04", "1 vCPU", "2 GB", "20 GB SSD", "running", "us-east-1", "60d 11h", "Security"},
	{"ci-runner-01", "10.0.7.40", "Ubuntu 24.04", "8 vCPU", "16 GB", "500 GB SSD", "running", "us-west-2", "3d 7h", "DevOps"},
	{"ci-runner-02", "10.0.7.41", "Ubuntu 24.04", "8 vCPU", "16 GB", "500 GB SSD", "running", "us-west-2", "3d 7h", "DevOps"},
}

type colDef struct {
	key   string
	title string
	field func(server) string
}

var colDefs = []colDef{
	{colName, "Name", func(s server) string { return s.name }},
	{colIP, "IP Address", func(s server) string { return s.ip }},
	{colOS, "OS", func(s server) string { return s.os }},
	{colCPU, "CPU", func(s server) string { return s.cpu }},
	{colRAM, "RAM", func(s server) string { return s.ram }},
	{colDisk, "Disk", func(s server) string { return s.disk }},
	{colStatus, "Status", func(s server) string { return s.status }},
	{colRegion, "Region", func(s server) string { return s.region }},
	{colUptime, "Uptime", func(s server) string { return s.uptime }},
	{colOwner, "Owner", func(s server) string { return s.owner }},
}

func maxLen(header string, field func(server) string) int {
	n := len(header)
	for _, s := range servers {
		if l := len(field(s)); l > n {
			n = l
		}
	}
	return n
}

type model struct {
	table teagrid.Model
}

func newModel() model {
	columns := make([]teagrid.Column, len(colDefs))
	for i, cd := range colDefs {
		columns[i] = teagrid.NewColumn(cd.key, cd.title, maxLen(cd.title, cd.field))
	}

	rows := make([]teagrid.Row, len(servers))
	for i, s := range servers {
		rows[i] = teagrid.NewRow(teagrid.RowData{
			colName:   s.name,
			colIP:     s.ip,
			colOS:     s.os,
			colCPU:    s.cpu,
			colRAM:    s.ram,
			colDisk:   s.disk,
			colStatus: s.status,
			colRegion: s.region,
			colUptime: s.uptime,
			colOwner:  s.owner,
		})
	}

	return model{
		table: teagrid.
			New(columns).
			WithRows(rows).
			WithHorizontalFreezeColumnCount(1).
			WithStaticFooter("Scroll: shift+left/right | Navigate: up/down | Quit: q").
			Focused(true),
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
		m.table = m.table.SetSize(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var body strings.Builder

	body.WriteString("Server inventory with horizontal scrolling and frozen Name column.\n")
	body.WriteString("Press shift+left/right to scroll, q or ctrl+c to quit.\n\n")
	body.WriteString(m.table.View().Content)

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
