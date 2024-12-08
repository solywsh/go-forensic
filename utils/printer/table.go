package printer

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"sync"
)

var DefaultTableBaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type (
	TableModelX struct {
		m       *tableModel
		p       *tea.Program
		startWg *sync.WaitGroup
	}

	tableModel struct {
		columns        []table.Column
		rows           []table.Row
		tableStyle     table.Styles
		table          table.Model
		height         int
		width          int
		tableBaseStyle lipgloss.Style
		ctx            context.Context
		cancel         context.CancelFunc
		selectedRow    chan table.Row
	}
)

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			select {
			// non-blocking data transmission
			case m.selectedRow <- m.table.SelectedRow():
			default:
			}
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return m.tableBaseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func NewTable(ctx context.Context) *TableModelX {
	_ctx, cancel := context.WithCancel(ctx)
	return &TableModelX{
		m: &tableModel{
			tableBaseStyle: DefaultTableBaseStyle,
			tableStyle: func() table.Styles {
				s := table.DefaultStyles()
				s.Header = s.Header.
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("240")).
					BorderBottom(true).
					Bold(false)
				s.Selected = s.Selected.
					Foreground(lipgloss.Color("229")).
					Background(lipgloss.Color("57")).
					Bold(false)
				return s
			}(),
			height:      10,
			width:       0,
			ctx:         _ctx,
			cancel:      cancel,
			selectedRow: make(chan table.Row, 10),
		},
		startWg: &sync.WaitGroup{},
	}
}

func (t *TableModelX) SetColumns(columns []table.Column) *TableModelX {
	t.m.columns = columns
	return t
}

func (t *TableModelX) SetRows(rows []table.Row) *TableModelX {
	t.m.rows = rows
	return t
}

func (t *TableModelX) SetHeight(height int) *TableModelX {
	t.m.height = height
	return t
}

func (t *TableModelX) SetWidth(width int) *TableModelX {
	t.m.width = width
	return t
}

func (t *TableModelX) SetTableStyle(style table.Styles) *TableModelX {
	t.m.tableStyle = style
	return t
}

func (t *TableModelX) Run() {
	t.startWg.Add(1)
	go func() {
		defer t.m.cancel()
		defer close(t.m.selectedRow)
		tb := table.New(
			table.WithColumns(t.m.columns),
			table.WithRows(t.m.rows),
			table.WithFocused(true),
			table.WithHeight(t.m.height),
			table.WithWidth(t.m.width))
		tb.SetStyles(t.m.tableStyle)
		t.m.table = tb
		t.p = tea.NewProgram(t.m)
		t.startWg.Done()
		if _, err := t.p.Run(); err != nil {
			fmt.Println(err)
		}
	}()
}

func (t *TableModelX) Quit() {
	t.startWg.Wait()
	if t.p != nil {
		t.p.Quit()
	}
}

func (t *TableModelX) Wait() {
	select {
	case <-t.m.ctx.Done():
		t.Quit()
	}
}

func (t *TableModelX) Ctx() context.Context {
	return t.m.ctx
}

func (t *TableModelX) SelectedRowChan() chan table.Row {
	return t.m.selectedRow
}

func (t *TableModelX) SelectedRow() (table.Row, bool) {
	select {
	case row, ok := <-t.m.selectedRow:
		return row, ok
	default:
		return nil, false
	}
}
