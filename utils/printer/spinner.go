package printer

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sync"
)

var (
	defaultTextStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	defaultSpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type SpinnerX struct {
	m  *SpinnerModel
	p  *tea.Program
	wg *sync.WaitGroup
}

func NewSpinnerX() SpinnerX {
	return SpinnerX{
		m:  newSpinnerModel(),
		wg: &sync.WaitGroup{},
	}
}

func (s SpinnerX) Run() {
	s.wg.Add(1)
	go func() {
		s.p = tea.NewProgram(s.m)
		s.wg.Done()
		_, err := s.p.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}

func (s SpinnerX) Msg(msg string) SpinnerX {
	s.m.msg = msg
	s.wg.Wait()
	if s.p != nil {
		s.p.Send(updateMessage{text: msg})
	}
	return s
}

func (s SpinnerX) Quit() {
	s.wg.Wait()
	if s.p != nil {
		s.p.Quit()
	}
}

type SpinnerModel struct {
	msg          string
	textStyle    lipgloss.Style
	spinnerStyle lipgloss.Style
	spinner      spinner.Model
}

type updateMessage struct {
	text string
}

func newSpinnerModel() *SpinnerModel {
	s := SpinnerModel{
		spinner:   spinner.New(),
		textStyle: defaultTextStyle,
	}
	s.spinner.Style = defaultSpinnerStyle
	return &s
}

func (m *SpinnerModel) SetSpinner(spinner spinner.Spinner) *SpinnerModel {
	if m == nil {
		return nil
	}
	m.spinner.Spinner = spinner
	return m
}

func (m *SpinnerModel) SetSpinnerStyle(style lipgloss.Style) *SpinnerModel {
	if m == nil {
		return nil
	}
	m.spinner.Style = style
	return m
}

func (m *SpinnerModel) SetTextStyle(style lipgloss.Style) *SpinnerModel {
	if m == nil {
		return nil
	}
	m.textStyle = style
	return m
}

func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case updateMessage:
		m.msg = msg.text
	}
	return m, nil
}

func (m *SpinnerModel) View() string {
	return fmt.Sprintf("%s %s\n", (m.spinner).View(), m.textStyle.Render(m.msg))
}
