// Copyright 2025 Christopher O'Connell
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package views

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/christopher/mcl/pkg/container"
	"github.com/christopher/mcl/pkg/tui/style"
)

// HomeModel is the main container list view
type HomeModel struct {
	table         table.Model
	width         int
	height        int
	animState     int
	containers    []container.Info
	daemonRunning bool
}

// NewHomeModel creates a new home view
func NewHomeModel(containers []container.Info, daemonRunning bool) *HomeModel {
	columns := []table.Column{
		{Title: "NAME", Width: 25},
		{Title: "STATUS", Width: 14},
		{Title: "BRANCH", Width: 25},
		{Title: "GIT", Width: 10},
		{Title: "ACTIVITY", Width: 12},
		{Title: "AUTH", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Custom styles with Ocean Tide colors
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(style.PurpleHaze).
		BorderBottom(true).
		Bold(true).
		Foreground(style.OceanTide)

	s.Selected = s.Selected.
		Foreground(style.GhostWhite).
		Background(lipgloss.Color("237")).
		Bold(false)

	t.SetStyles(s)

	h := &HomeModel{
		table:         t,
		containers:    containers,
		daemonRunning: daemonRunning,
	}

	h.updateTableRows()
	return h
}

// Init initializes the home view
func (h *HomeModel) Init() tea.Cmd {
	return nil
}

// Update handles input and state changes
func (h *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return h, tea.Quit
		case "enter":
			// Get selected container
			if len(h.containers) > 0 {
				selectedIdx := h.table.Cursor()
				if selectedIdx >= 0 && selectedIdx < len(h.containers) {
					selected := h.containers[selectedIdx]
					// Return a message to signal connection request
					return h, func() tea.Msg {
						return ConnectRequestMsg{ContainerName: selected.Name}
					}
				}
			}
			return h, nil
		case "a":
			// Show actions menu for selected container
			if len(h.containers) > 0 {
				selectedIdx := h.table.Cursor()
				if selectedIdx >= 0 && selectedIdx < len(h.containers) {
					selected := h.containers[selectedIdx]
					return h, func() tea.Msg {
						return ShowActionsMenuMsg{Container: selected}
					}
				}
			}
			return h, nil
		case "up", "k":
			h.table, cmd = h.table.Update(msg)
			return h, cmd
		case "down", "j":
			h.table, cmd = h.table.Update(msg)
			return h, cmd
		}
	}

	h.table, cmd = h.table.Update(msg)
	return h, cmd
}

// ConnectRequestMsg signals that the user wants to connect to a container
type ConnectRequestMsg struct {
	ContainerName string
}

// ShowActionsMenuMsg signals to show the actions menu for a container
type ShowActionsMenuMsg struct {
	Container container.Info
}

// View renders the home view
func (h *HomeModel) View() string {
	// Container table
	tableView := h.table.View()

	// Center the table horizontally
	return lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Top,
		tableView,
	)
}

// SetSize updates the view dimensions
func (h *HomeModel) SetSize(width, height int) {
	h.width = width
	h.height = height

	// Adjust table height to fill screen
	// Title (1) + empty (1) + empty (1) + help bar (1) = 4 lines overhead
	tableHeight := height - 4
	if tableHeight < 5 {
		tableHeight = 5
	}
	// Don't limit by container count - let table scroll if needed
	h.table.SetHeight(tableHeight)

	// Set table width
	h.table.SetWidth(width)
}

// SetAnimationState updates the animation state for pulsing indicators
func (h *HomeModel) SetAnimationState(state int) {
	h.animState = state
}

// RefreshContainers updates the container list
func (h *HomeModel) RefreshContainers(containers []container.Info, daemonRunning bool) {
	h.containers = containers
	h.daemonRunning = daemonRunning
	h.updateTableRows()
}

// updateTableRows converts container data to table rows
func (h *HomeModel) updateTableRows() {
	rows := make([]table.Row, 0, len(h.containers))

	for _, c := range h.containers {
		row := table.Row{
			h.formatName(c),
			h.formatStatus(c),
			h.formatBranch(c),
			h.formatGit(c),
			h.formatActivity(c),
			h.formatAuth(c),
		}
		rows = append(rows, row)
	}

	h.table.SetRows(rows)
}

// formatName returns the container short name
func (h *HomeModel) formatName(c container.Info) string {
	return c.ShortName
}

// formatStatus returns the status indicator
func (h *HomeModel) formatStatus(c container.Info) string {
	switch c.Status {
	case "running":
		return "● Running"
	case "exited":
		return "■ Stopped"
	default:
		return "? " + c.Status
	}
}

// formatBranch returns the branch name
func (h *HomeModel) formatBranch(c container.Info) string {
	if c.Branch == "" {
		return "—"
	}
	return c.Branch
}

// formatGit returns git status
func (h *HomeModel) formatGit(c container.Info) string {
	if c.GitStatus == "" {
		return "—"
	}
	return c.GitStatus
}

// formatActivity returns time since last activity
func (h *HomeModel) formatActivity(c container.Info) string {
	if c.LastActivity == "" {
		return "—"
	}
	return c.LastActivity
}

// formatAuth returns authentication status
func (h *HomeModel) formatAuth(c container.Info) string {
	if c.AuthStatus == "" {
		return "—"
	}
	return c.AuthStatus
}

// GetContainers returns the current container list for caching
func (h *HomeModel) GetContainers() []container.Info {
	return h.containers
}

// GetCursor returns the current cursor position for caching
func (h *HomeModel) GetCursor() int {
	return h.table.Cursor()
}

// SetCursor sets the cursor position (used when restoring from cache)
func (h *HomeModel) SetCursor(pos int) {
	if pos >= 0 && pos < len(h.containers) {
		h.table.SetCursor(pos)
	}
}
