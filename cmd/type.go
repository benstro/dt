package cmd

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/benstro/dt/internal/typer"
)

const (
	defaultWordCount = 25
	timeModeBuffer   = 300 // words pre-generated for time-attack mode
	lineWidth        = 65  // max chars per visual line
)

// ── phase ────────────────────────────────────────────────────────────────────

type phase int

const (
	phaseReady    phase = iota
	phasePlaying
	phaseFinished
)

// ── model ────────────────────────────────────────────────────────────────────

type tickMsg time.Time

type typeModel struct {
	target    []rune
	input     []rune
	phase     phase
	wordCount int
	timeLimit int // seconds; 0 = word-count mode
	startTime time.Time
	elapsed   time.Duration
	result    typer.Result
	width     int
}

func newTypeModel(wordCount, timeLimit int) typeModel {
	n := wordCount
	if timeLimit > 0 {
		n = timeModeBuffer
	}
	return typeModel{
		target:    []rune(strings.Join(typer.Sample(n), " ")),
		wordCount: wordCount,
		timeLimit: timeLimit,
		width:     80,
	}
}

func (m typeModel) restart() typeModel {
	return newTypeModel(m.wordCount, m.timeLimit)
}

// ── bubbletea ────────────────────────────────────────────────────────────────

func (m typeModel) Init() tea.Cmd { return nil }

func (m typeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tickMsg:
		if m.phase != phasePlaying {
			return m, nil
		}
		m.elapsed = time.Since(m.startTime)
		if m.timeLimit > 0 && m.elapsed >= time.Duration(m.timeLimit)*time.Second {
			m.elapsed = time.Duration(m.timeLimit) * time.Second
			m.phase = phaseFinished
			m.result = typer.Score(string(m.target), string(m.input), m.elapsed)
			return m, nil
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch m.phase {
		case phaseReady:
			return m.updateReady(msg)
		case phasePlaying:
			return m.updatePlaying(msg)
		case phaseFinished:
			return m.updateFinished(msg)
		}
	}
	return m, nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m typeModel) updateReady(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyTab:
		return m.restart(), nil
	case tea.KeyRunes:
		m.phase = phasePlaying
		m.startTime = time.Now()
		m = m.press(k.Runes[0])
		return m, tickCmd()
	}
	return m, nil
}

func (m typeModel) updatePlaying(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyTab:
		return m.restart(), nil
	case tea.KeyBackspace:
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
		return m, nil
	case tea.KeySpace:
		m.elapsed = time.Since(m.startTime)
		return m.pressAndCheck(' ')
	case tea.KeyRunes:
		m.elapsed = time.Since(m.startTime)
		return m.pressAndCheck(k.Runes[0])
	}
	return m, nil
}

func (m typeModel) pressAndCheck(r rune) (typeModel, tea.Cmd) {
	m = m.press(r)
	if m.timeLimit == 0 && len(m.input) >= len(m.target) {
		m.phase = phaseFinished
		m.result = typer.Score(string(m.target), string(m.input), m.elapsed)
		return m, nil
	}
	return m, nil
}

func (m typeModel) press(r rune) typeModel {
	if m.timeLimit == 0 && len(m.input) >= len(m.target) {
		return m
	}
	m.input = append(m.input, r)
	return m
}

func (m typeModel) updateFinished(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit
	case tea.KeyTab:
		return m.restart(), nil
	}
	return m, nil
}

// ── styles ───────────────────────────────────────────────────────────────────

var (
	sDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("246")) // light gray — readable on Catppuccin dark bg
	sCorrect = lipgloss.NewStyle().Foreground(lipgloss.Color("255")) // near-white — clearly typed
	sWrong   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	sCursor  = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Underline(true)
	sLabel   = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	sValue   = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	sHint    = lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
)

// ── view ─────────────────────────────────────────────────────────────────────

func (m typeModel) View() string {
	if m.phase == phaseFinished {
		return m.viewResults()
	}
	return m.viewGame()
}

func (m typeModel) viewGame() string {
	w := lineWidth
	if m.width-4 < w {
		w = m.width - 4
	}
	lines := wrapRunes(m.target, w)

	// Find which visual line the cursor is on.
	cursorPos := len(m.input)
	offset := 0
	cursorLine := len(lines) - 1
	for i, l := range lines {
		if cursorPos < offset+len(l) {
			cursorLine = i
			break
		}
		offset += len(l)
	}

	// Show a 3-line window: line before, cursor line, line after.
	start := cursorLine - 1
	if start < 0 {
		start = 0
	}
	end := start + 3
	if end > len(lines) {
		end = len(lines)
	}

	// Character offset at the start of the window.
	charOffset := 0
	for i := 0; i < start; i++ {
		charOffset += len(lines[i])
	}

	var sb strings.Builder
	sb.WriteString("\n")
	for li := start; li < end; li++ {
		sb.WriteString("  ")
		for ci, r := range lines[li] {
			gi := charOffset + ci
			switch {
			case gi == cursorPos:
				sb.WriteString(sCursor.Render(string(r)))
			case gi < len(m.input):
				if m.input[gi] == r {
					sb.WriteString(sCorrect.Render(string(r)))
				} else {
					sb.WriteString(sWrong.Render(string(r)))
				}
			default:
				sb.WriteString(sDim.Render(string(r)))
			}
		}
		sb.WriteString("\n")
		charOffset += len(lines[li])
	}
	// Pad to 3 lines to prevent layout shift.
	for i := end - start; i < 3; i++ {
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Stats row.
	if m.phase == phasePlaying {
		r := typer.Score(string(m.target), string(m.input), m.elapsed)
		var timer string
		if m.timeLimit > 0 {
			rem := time.Duration(m.timeLimit)*time.Second - m.elapsed
			if rem < 0 {
				rem = 0
			}
			timer = fmt.Sprintf("%d", int(rem.Seconds()))
		} else {
			s := int(m.elapsed.Seconds())
			timer = fmt.Sprintf("%d:%02d", s/60, s%60)
		}
		sb.WriteString(fmt.Sprintf("  %s %s   %s %s%%   %s %s\n",
			sLabel.Render("wpm"), sValue.Render(fmt.Sprintf("%d", r.NetWPM)),
			sLabel.Render("acc"), sValue.Render(fmt.Sprintf("%.0f", r.Accuracy)),
			sLabel.Render("time"), sValue.Render(timer),
		))
	} else {
		sb.WriteString(sHint.Render("  start typing...") + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(sHint.Render("  tab: new game  ·  esc: quit") + "\n")
	return sb.String()
}

func (m typeModel) viewResults() string {
	r := m.result
	s := int(r.Elapsed.Seconds())

	var sb strings.Builder
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("  %s  %s\n", sLabel.Render("wpm "), sValue.Render(fmt.Sprintf("%d", r.NetWPM))))
	sb.WriteString(fmt.Sprintf("  %s  %s%%\n", sLabel.Render("acc "), sValue.Render(fmt.Sprintf("%.0f", r.Accuracy))))
	sb.WriteString(fmt.Sprintf("  %s  %s\n", sLabel.Render("raw "), sValue.Render(fmt.Sprintf("%d", r.RawWPM))))
	sb.WriteString(fmt.Sprintf("  %s  %s\n", sLabel.Render("time"), sValue.Render(fmt.Sprintf("%d:%02d", s/60, s%60))))
	sb.WriteString("\n")
	sb.WriteString(sHint.Render("  tab: new game  ·  esc: quit") + "\n")
	return sb.String()
}

// wrapRunes splits target into lines of at most width runes, breaking at spaces.
// The breaking space is included at the end of its line so character indices
// stay contiguous and the user must type it.
func wrapRunes(target []rune, width int) [][]rune {
	var lines [][]rune
	i := 0
	for i < len(target) {
		if i+width >= len(target) {
			lines = append(lines, target[i:])
			break
		}
		// Find the last space within [i, i+width].
		j := i + width
		for j > i && target[j] != ' ' {
			j--
		}
		if j == i {
			j = i + width // hard break — no space found
		} else {
			j++ // include the space on this line
		}
		lines = append(lines, target[i:j])
		i = j
	}
	return lines
}

// ── cobra command ─────────────────────────────────────────────────────────────

var typeCmd = &cobra.Command{
	Use:   "type",
	Short: "Typing speed test",
	RunE: func(cmd *cobra.Command, args []string) error {
		words, _ := cmd.Flags().GetInt("words")
		seconds, _ := cmd.Flags().GetInt("time")

		m := newTypeModel(words, seconds)
		p := tea.NewProgram(m, tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func init() {
	typeCmd.Flags().IntP("words", "w", defaultWordCount, "number of words (word-count mode)")
	typeCmd.Flags().IntP("time", "t", 0, "time limit in seconds (time-attack mode, e.g. -t 60)")
	typeCmd.MarkFlagsMutuallyExclusive("words", "time")
	rootCmd.AddCommand(typeCmd)
}
