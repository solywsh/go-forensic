package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver
)

type SqliteX struct {
	dbPath   string
	keywords []string
}

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
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		// get column information for each table
		columnRows, err := db.Query(fmt.Sprintf("PRAGMA table_info(\"%s\");", tableName))
		if err != nil {
			// 如果获取列信息失败，跳过该表，继续检查下一个表
			//fmt.Printf("发生错误: 获取表 %s 列信息失败: %v\n", tableName, err)
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
					fmt.Printf("found keyword '%s' in column '%s' of table '%s'\n", tableName, columnName, keyword)
				}
				rowsInColumn.Close()
			}
		}
		columnRows.Close()
	}
	return nil
}
