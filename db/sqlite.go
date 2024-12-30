package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
	"github.com/solywsh/go-forensic/utils"
	"github.com/solywsh/go-forensic/utils/osx"
	"github.com/solywsh/go-forensic/utils/printer"
)

type (
	SqliteX struct {
		dbPath   string
		keywords []string

		tableMaxWidth  int
		columnMaxWidth int
		keyMaxWidth    int
		valueMaxWidth  int // TODO show value
		ignoreErr      bool

		tableHeight int
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
	return &SqliteX{
		tableMaxWidth:  6,
		columnMaxWidth: 6,
		keyMaxWidth:    6,
		valueMaxWidth:  6,
		ignoreErr:      false,
		tableHeight:    10,
	}
}

func (t *SqliteX) SetDbPath(dbPath string) *SqliteX {
	t.dbPath = dbPath
	return t
}

func (t *SqliteX) SetIgnoreErr(ignoreErr bool) *SqliteX {
	t.ignoreErr = ignoreErr
	return t
}

func (t *SqliteX) SetTableHeight(tableHeight int) *SqliteX {
	t.tableHeight = tableHeight
	return t
}

func (t *SqliteX) Error(msg interface{}, keyvals ...interface{}) {
	if t.ignoreErr {
		return
	}
	log.Error(msg, keyvals...)
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
			t.Error("failed to scan table name", "error", err)
			continue
		}
		// get column information for each table
		columnRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(\"%s\");", tableName))
		if err != nil {
			t.Error("failed to get column information", "table", tableName, "error", err)
			continue
		}
		var columnName string
		for columnRows.Next() {
			var cid int
			var cType, notnull, pk string
			var defaultValue sql.NullString
			// scan column information
			if err := columnRows.Scan(&cid, &columnName, &cType, &notnull, &defaultValue, &pk); err != nil {
				t.Error("failed to scan column information", "table", tableName, "error", err)
				continue
			}
			// search for keywords in each column
			query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIKE ?;", tableName, columnName)
			for _, keyword := range keywords {
				rowsInColumn, err := db.Query(query, "%"+keyword+"%")
				if err != nil {
					t.Error("failed to query column %s of table %s: %v", tableName, columnName, err)
					continue
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
					t.tableMaxWidth = utils.Max(t.tableMaxWidth, len(tableName))
					t.columnMaxWidth = utils.Max(t.columnMaxWidth, len(columnName))
					t.keyMaxWidth = utils.Max(t.keyMaxWidth, len(keyword))
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
		{Title: "Table", Width: utils.Min(50, utils.Max(t.tableMaxWidth, len("Table")))},
		{Title: "Column", Width: utils.Min(50, utils.Max(t.columnMaxWidth, len("Column")))},
		{Title: "Keywords", Width: utils.Min(50, utils.Max(t.keyMaxWidth, len("Keywords")))},
	}
	rows := make([]table.Row, 0, len(res))
	for _, r := range res {
		rows = append(rows, table.Row{r.Table, r.Column, r.Keywords})
	}
	tb := printer.NewTable(context.Background())
	if t.tableHeight > 0 {
		tb.SetHeight(t.tableHeight)
	}
	tb.SetColumns(columns).SetRows(rows)
	tb.Run()
	tb.Wait()
	return nil
}
