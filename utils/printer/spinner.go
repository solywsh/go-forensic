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

type (
	SpinnerX struct {
		m  *SpinnerModel
		p  *tea.Program
		wg *sync.WaitGroup
	}

	SpinnerModel struct {
		msg          string
		textStyle    lipgloss.Style
		spinnerStyle lipgloss.Style
		spinner      spinner.Model
	}

	updateMessage struct {
		text string
	}
	updateSpinner struct {
		spinner spinner.Spinner
	}
	updateSpinnerStyle struct {
		style lipgloss.Style
	}
	updateTextStyle struct {
		style lipgloss.Style
	}
)

func NewSpinnerX() *SpinnerX {
	return &SpinnerX{
		m:  initSpinnerModel(),
		wg: &sync.WaitGroup{},
	}
}

func initSpinnerModel() *SpinnerModel {
	s := SpinnerModel{
		spinner:   spinner.New(),
		textStyle: defaultTextStyle,
	}
	s.spinner.Style = defaultSpinnerStyle
	return &s
}

func (s *SpinnerX) Run() {
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

func (s *SpinnerX) Msg(msg string) *SpinnerX {
	s.m.msg = msg
	s.wg.Wait()
	if s.p != nil {
		s.p.Send(updateMessage{text: msg})
	}
	return s
}

func (s *SpinnerX) Quit() {
	s.wg.Wait()
	if s.p != nil {
		s.p.Quit()
	}
}

func (s *SpinnerX) SetSpinner(spinner spinner.Spinner) *SpinnerX {
	if s == nil {
		return nil
	}
	s.wg.Wait()
	if s.p != nil {
		s.p.Send(updateSpinner{spinner: spinner})
	}
	return s
}

func (s *SpinnerX) SetSpinnerStyle(style lipgloss.Style) *SpinnerX {
	if s == nil {
		return nil
	}
	s.wg.Wait()
	if s.p != nil {
		s.p.Send(updateSpinnerStyle{style: style})
	}
	return s
}

func (s *SpinnerX) SetTextStyle(style lipgloss.Style) *SpinnerX {
	if s == nil {
		return nil
	}
	s.wg.Wait()
	if s.p != nil {
		s.p.Send(updateTextStyle{style: style})
	}
	return s
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
	case updateSpinner:
		m.spinner.Spinner = msg.spinner
	case updateSpinnerStyle:
		m.spinner.Style = msg.style
	case updateTextStyle:
		m.textStyle = msg.style
	}
	return m, nil
}

func (m *SpinnerModel) View() string {
	return fmt.Sprintf("%s %s\n", (m.spinner).View(), m.textStyle.Render(m.msg))
}
