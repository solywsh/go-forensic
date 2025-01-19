package printer

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/solywsh/go-forensic/utils"
	"strings"
	"sync"
)

type (
	PaperX struct {
		m       *paperModel
		p       *tea.Program
		startWg *sync.WaitGroup
	}
	paperModel struct {
		content      string
		header       string
		loading      string
		footer       string // TODO Customize footer content
		border       string
		showProgress bool
		ready        bool
		viewport     viewport.Model
		// You generally won't need this unless you're processing stuff with
		// complicated ANSI escape sequences. Turn it on if you notice flickering.
		//
		// Also keep in mind that high performance rendering only works for programs
		// that use the full size of the terminal. We're enabling that below with
		// tea.EnterAltScreen().
		useHighPerformanceRenderer bool
		titleStyle                 lipgloss.Style
		infoStyle                  lipgloss.Style
		ctx                        context.Context
		cancel                     context.CancelFunc
	}
)

var (
	DefaultPaperTitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	DefaultPaperInfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return DefaultPaperTitleStyle.BorderStyle(b)
	}()
)

func (m paperModel) Init() tea.Cmd {
	return nil
}

func (m paperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd     tea.Cmd
		cmdList []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = m.useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		if m.useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmdList = append(cmdList, viewport.Sync(m.viewport))
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmdList = append(cmdList, cmd)

	return m, tea.Batch(cmdList...)
}

func (m paperModel) View() string {
	if !m.ready {
		if m.loading != "" {
			return "\n " + m.loading
		} else {
			return "\n  Initializing..."
		}
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())

}

func (m paperModel) headerView() string {
	if m.header != "" {
		title := m.titleStyle.Render(m.header)
		line := strings.Repeat(m.border, utils.Max(0, m.viewport.Width-lipgloss.Width(title)))
		return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	} else {
		return lipgloss.JoinHorizontal(lipgloss.Center, strings.Repeat(m.border, m.viewport.Width))
	}
}

func (m paperModel) footerView() string {
	if m.showProgress {
		info := m.infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
		line := strings.Repeat(m.border, utils.Max(0, m.viewport.Width-lipgloss.Width(info)))
		return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	} else {
		return lipgloss.JoinHorizontal(lipgloss.Center, strings.Repeat(m.border, m.viewport.Width))
	}
}

func (p *PaperX) SetTitle(title string) *PaperX {
	p.m.header = title
	return p
}

func (p *PaperX) SetContent(content string) *PaperX {
	p.m.content = content
	return p
}

func (p *PaperX) SetLoading(loading string) *PaperX {
	p.m.loading = loading
	return p
}

func (p *PaperX) EnableProgress(enable bool) *PaperX {
	p.m.showProgress = enable
	return p
}

func (p *PaperX) EnableHighPerformanceRendering(enable bool) *PaperX {
	p.m.useHighPerformanceRenderer = enable
	return p
}

func (p *PaperX) SetBorder(border string) *PaperX {
	p.m.border = border
	return p
}

func (p *PaperX) SetTitleStyle(style lipgloss.Style) *PaperX {
	p.m.titleStyle = style
	return p
}

func (p *PaperX) SetInfoStyle(style lipgloss.Style) *PaperX {
	p.m.infoStyle = style
	return p
}

func NewPaper(ctx context.Context) *PaperX {
	_ctx, cancel := context.WithCancel(ctx)
	return &PaperX{m: &paperModel{
		showProgress:               true,
		border:                     "-",
		loading:                    "load content...",
		useHighPerformanceRenderer: false,
		titleStyle:                 DefaultPaperTitleStyle,
		infoStyle:                  DefaultPaperInfoStyle,
		ctx:                        _ctx,
		cancel:                     cancel,
	},
		startWg: &sync.WaitGroup{},
	}
}

func (p *PaperX) Run() {
	p.startWg.Add(1)
	go func() {
		defer p.m.cancel()
		p.p = tea.NewProgram(
			p.m,
			tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
			tea.WithMouseCellMotion(), // turn on mouse support, so we can track the mouse wheel
		)
		p.startWg.Done()
		_, err := p.p.Run()
		if err != nil {
			log.Error(err)
			return
		}
	}()
}

func (p *PaperX) Quit() {
	p.startWg.Wait()
	if p.p != nil {
		p.p.Quit()
	}
}

func (p *PaperX) Wait() {
	select {
	case <-p.m.ctx.Done():
		p.Quit()
	}
}

func (p *PaperX) Ctx() context.Context {
	return p.m.ctx
}
