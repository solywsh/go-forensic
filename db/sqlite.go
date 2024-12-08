package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
)

type (
	SqliteX struct {
		dbPath   string
		keywords []string

		tableMaxLength  int
		columnMaxLength int
		keyMaxLength    int
		valueMaxLength  int // TODO show value
	}
	SearchResult struct {
		Path     string
		Table    string
		Column   string
		Keywords string
		Key      string // other database
		Value    string
	}
)

func NewSqlite() *SqliteX {
	return &SqliteX{}
}

func (t *SqliteX) SetDbPath(dbPath string) *SqliteX {
	t.dbPath = dbPath
	return t
}

func (t *SqliteX) SearchByKeywords(keywords ...string) error {
	if t.dbPath == "" {
		return fmt.Errorf("db path is empty")
	}
	db, err := sql.Open("sqlite3", t.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		return err
	}
	defer rows.Close()
	var tableName string
	var res []SearchResult
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		// get column information for each table
		columnRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(\"%s\");", tableName))
		if err != nil {
			log.Error("failed to get column information", "table", tableName, "error", err)
			continue
		}
		var columnName string
		for columnRows.Next() {
			var cid int
			var cType, notnull, pk string
			var defaultValue sql.NullString
			// scan column information
			if err := columnRows.Scan(&cid, &columnName, &cType, &notnull, &defaultValue, &pk); err != nil {
				return err
			}
			// search for keywords in each column
			query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIKE ?;", tableName, columnName)
			for _, keyword := range keywords {
				rowsInColumn, err := db.Query(query, "%"+keyword+"%")
				if err != nil {
					return fmt.Errorf("failed to query column %s of table %s: %v", tableName, columnName, err)
				}
				// If a record containing the keyword is found in this column, output the information.
				if rowsInColumn.Next() {
					if !osx.IsTTY() {
						log.Info("found keyword", "table", tableName, "column", columnName, "keyword", keyword)
					}
					res = append(res, SearchResult{
						Table:    tableName,
						Column:   columnName,
						Keywords: keyword,
					})
					t.tableMaxLength = max(t.tableMaxLength, len(tableName))
					t.columnMaxLength = max(t.columnMaxLength, len(columnName))
					t.keyMaxLength = max(t.keyMaxLength, len(keyword))
				}
				rowsInColumn.Close()
			}
		}
		columnRows.Close()
	}
	err = t.showSearchResult(res)
	if err != nil {
		return err
	}
	return nil
}

func (t *SqliteX) showSearchResult(res []SearchResult) error {
	if len(res) == 0 {
		log.Info("no search result")
		return nil
	}
	if !osx.IsTTY() {
		return nil
	}
	columns := []table.Column{
		{Title: "Table", Width: min(50, t.tableMaxLength)},
		{Title: "Column", Width: min(50, t.columnMaxLength)},
		{Title: "Keywords", Width: min(50, t.keyMaxLength)},
	}
	rows := make([]table.Row, 0, len(res))
	for _, r := range res {
		rows = append(rows, table.Row{r.Table, r.Column, r.Keywords})
	}
	tb := printer.NewTable(context.Background())
	tb.SetColumns(columns).SetRows(rows)
	tb.Run()
	tb.Wait()
	return nil
}
