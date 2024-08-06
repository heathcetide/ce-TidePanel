package router

import (
	"fmt"
	"github.com/heathcetide/ce-TidePanel/config"

	"os"
	"testing"
)

func TestGetTableData(t *testing.T) {
	db, _ := config.GetDB()
	tableName := "users"
	sqlString := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
	rows, err := db.Raw(sqlString).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var column ColumnInfo
		if err := rows.Scan(
			&column.Field,
			&column.Type,
			&column.Null,
			&column.Key,
			&column.Default,
			&column.Extra,
		); err != nil {
			return
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return
	}
	structStr := generateStruct(columns, modifyString(tableName))
	// 文件名
	filename := fmt.Sprintf("../model/%s.go", tableName)
	// 使用 os.OpenFile 打开或创建文件
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error opening/creating file:", err)
		return
	}
	defer file.Close()
	content := fmt.Sprintf("package main\n\n import (\"time\")\n\n %s", structStr)
	// 写入文件
	_, err = file.WriteString(content)
}
