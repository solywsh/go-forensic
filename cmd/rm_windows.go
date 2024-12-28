//go:build windows
// +build windows

package cmd

import (
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/spf13/cobra"
)

var (
	rmPath string
	rmCmd  = &cobra.Command{
		Use:   "rm",
		Short: "Delete files under the specified path (Windows only)",
		Run: func(cmd *cobra.Command, args []string) {
			if rmPath == "" {
				return
			}
			osx.RemoveAll(rmPath)
		},
	}
)

func init() {
	rmCmd.Flags().StringVarP(&rmPath, "path", "p", "", "the path to the file/directory to be deleted")
	rootCmd.AddCommand(rmCmd)
}
