package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"trello_cli/trello"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

type (
	errMsg error
)

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter your Trello API Key"
	inputs[0].Focus()
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Enter your Trello API Token"
	inputs[1].Width = 50

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs)-1 {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		`Trello CLI Configuration

%s
%s

%s
`,
		inputStyle.Width(60).Render("API Key: "+m.inputs[0].View()),
		inputStyle.Width(60).Render("API Token: "+m.inputs[1].View()),
		continueStyle.Render("Press Enter to continue"),
	) + "\n"
}

func PromptForConfig() (string, string, string, string, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return "", "", "", "", err
	}

	if m, ok := m.(model); ok {
		return m.inputs[0].Value(), m.inputs[1].Value(), "", "", nil
	}

	return "", "", "", "", fmt.Errorf("unexpected model type")
}

func PromptForOrganization(client *trello.Client) (string, error) {
	organizations, err := client.GetOrganizations()
	if err != nil {
		return "", fmt.Errorf("failed to fetch organizations: %w", err)
	}

	if len(organizations) == 0 {
		return "", fmt.Errorf("no organizations found")
	}

	fmt.Println("Available Workspaces:")
	for i, org := range organizations {
		fmt.Printf("%d. %s\n", i+1, org.Name)
	}

	var choice int
	fmt.Print("Select workspace (number): ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(organizations) {
		return "", fmt.Errorf("invalid choice")
	}

	return organizations[choice-1].ID, nil
}

func PromptForBoard(client *trello.Client, organizationID string) (string, error) {
	boards, err := client.GetBoards(organizationID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch boards: %w", err)
	}

	if len(boards) == 0 {
		return "", fmt.Errorf("no boards found in this workspace")
	}

	fmt.Println("Available Boards:")
	for i, board := range boards {
		fmt.Printf("%d. %s\n", i+1, board.Name)
	}

	var choice int
	fmt.Print("Select board (number): ")
	_, err = fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > len(boards) {
		return "", fmt.Errorf("invalid choice")
	}

	return boards[choice-1].ID, nil
}
