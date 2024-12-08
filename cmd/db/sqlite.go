package db

import (
	"github.com/charmbracelet/log"
	"github.com/solywsh/go-forensic/db"
	"github.com/spf13/cobra"
)

var (
	dbFilePath string
	keywords   []string
)

var (
	SqliteCmd = &cobra.Command{
		Use:   "sqlite",
		Short: "processing sqlite database related commands",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	SqliteSearchCmd = &cobra.Command{
		Use:   "search",
		Short: "search for information on keywords in the SQLite database",
		Run: func(cmd *cobra.Command, args []string) {
			err := db.NewSqlite().SetDbPath(dbFilePath).SearchByKeywords(keywords...)
			if err != nil {
				log.Error(err)
				return
			}
		},
	}
)

func init() {
	SqliteCmd.PersistentFlags().StringVarP(&dbFilePath, "file", "f", "", "the path to the SQLite database file")
	SqliteSearchCmd.Flags().StringSliceVarP(&keywords, "keywords", "k", nil, "the keyword to search for in the SQLite database")
	SqliteCmd.AddCommand(SqliteSearchCmd)
}
