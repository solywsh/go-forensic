package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/solywsh/go-forensic/utils/printer"
	"github.com/spf13/cobra"
	"runtime"
)

const Version = "0.0.1_beta"

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "version subcommand show go-forensic version info.",
		Run: func(cmd *cobra.Command, args []string) {
			vt := printer.NewText(Version).
				SetStyle(lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("#4c39c3")))
			st := printer.NewText(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)).
				SetStyle(lipgloss.
					NewStyle().
					Foreground(lipgloss.Color("#ffbca7")))
			fmt.Printf("go-forensic %s %s\n", vt.String(), st.String())
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
