package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const divisor = 4

type Model struct {
	loaded   bool
	focused  status
	lists    [3]list.Model
	err      error
	quitting bool
}

func New() *Model {
	return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
	var task Task
	var ok bool
	if task, ok = m.lists[m.focused].SelectedItem().(Task); !ok {
		return nil
	}
	m.lists[task.status].RemoveItem(m.lists[task.status].Index())
	task.Next()
	m.lists[task.status].InsertItem(-1, task)

	m.Update(nil)
	return nil
}

/*Go to next list */
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

/*Go to prev list */
func (m *Model) Prev() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height/2)
	defaultList.SetShowHelp(false)
	m.lists = [3]list.Model{defaultList, defaultList, defaultList}
	m.lists[0].Title = "To DO"
	m.lists[0].SetItems([]list.Item{
		Task{status: todo, title: "Test1", description: "Testing"},
		Task{status: todo, title: "Test2", description: "Testing"},
		Task{status: todo, title: "Test3", description: "Testing"},
	})
	m.lists[1].Title = "In Progress"
	m.lists[1].SetItems([]list.Item{
		Task{status: todo, title: "Test1", description: "Testing"},
		Task{status: todo, title: "Test2", description: "Testing"},
		Task{status: todo, title: "Test3", description: "Testing"},
	})
	m.lists[2].Title = "Done"
	m.lists[2].SetItems([]list.Item{
		Task{status: todo, title: "Test1", description: "Testing"},
		Task{status: todo, title: "Test2", description: "Testing"},
		Task{status: todo, title: "Test3", description: "Testing"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		}
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[0].View()
		inProgressView := m.lists[1].View()
		doneView := m.lists[2].View()

		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}

	} else {
		return "loading..."
	}
}

/*Styling */
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#2a9d8f"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)
