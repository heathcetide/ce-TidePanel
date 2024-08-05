package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"main/config"
	"strings"
)

// generateStruct 生成 Go 结构体
func generateStruct(columns []ColumnInfo, structStr string) string {
	var sb strings.Builder
	sb.WriteString("type " + structStr + " struct {\n")

	for _, col := range columns {
		goType := sqlTypeToGoType(*col.Type)
		sb.WriteString(fmt.Sprintf("\t%s %s `gorm:\"column:%s\"`\n", *col.Field, goType, *col.Field))
	}

	sb.WriteString("}\n")

	return sb.String()
}

// sqlTypeToGoType 将 SQL 类型转换为 Go 类型
func sqlTypeToGoType(sqlType string) string {
	switch sqlType {
	case "int", "int(11)", "tinyint", "smallint", "mediumint", "bigint":
		return "int"
	case "decimal", "numeric", "float", "real", "double", "double precision":
		return "float64"
	case "char", "varchar", "text", "tinytext", "mediumtext", "longtext":
		return "string"
	case "date", "datetime", "timestamp":
		return "time.Time"
	default:
		return "interface{}"
	}
}

func GetHtml(c *gin.Context) {
	db, err := config.GetDB()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "数据库连接失败",
		})
	}
	// 使用 db.Raw 执行 SQL 并获取结果
	rows, err := db.Raw("SHOW TABLES;").Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return
	}
	c.HTML(200, "test.html", tables)
}

type ColumnInfo struct {
	Field   *string
	Type    *string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}

func GetTableInfo(c *gin.Context) {
	db, err := config.GetDB()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "数据库连接失败",
		})
	}
	tableName := c.Param("tableName")
	columns := GetDBTable(tableName, db)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": columns,
	})
}
func modifyString(s string) string {
	if len(s) == 0 {
		return s
	}
	// Remove the last character
	s = s[:len(s)-1]

	// Capitalize the first character
	if len(s) > 0 {
		s = strings.ToUpper(s[:1]) + s[1:]
	}
	return s
}

// 辅助函数，用于获取指针指向的值或默认值
func getValue(s *string) string {
	if s == nil {
		return "NULL"
	}
	return *s
}

func GetTableData(c *gin.Context) {
	db, err := config.GetDB()
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "数据库连接失败",
		})
		return
	}
	tableName := c.Param("tableName")
	data := GetDBData(tableName, db)
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func GetDBTable(tableName string, db *gorm.DB) []ColumnInfo {
	sqlString := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
	rows, err := db.Raw(sqlString).Rows()
	if err != nil {
		return nil
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
			return nil
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	// 输出查询结果
	fmt.Println("Columns:")
	for _, col := range columns {
		fmt.Printf("Field: %s, Type: %s, Null: %s, Key: %s, Default: %s, Extra: %s\n",
			getValue(col.Field),
			getValue(col.Type),
			getValue(col.Null),
			getValue(col.Key),
			getValue(col.Default),
			getValue(col.Extra),
		)
	}
	return columns
}

// 获取指定表的数据
func GetDBData(tableName string, db *gorm.DB) []map[string]interface{} {
	sqlString := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Raw(sqlString).Rows()
	if err != nil {
		return nil
	}
	defer rows.Close()

	var data []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil
		}

		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				var v interface{}
				if err := json.Unmarshal(b, &v); err == nil {
					rowMap[colName] = v
				} else {
					rowMap[colName] = string(b)
				}
			} else {
				rowMap[colName] = val
			}
		}
		data = append(data, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return data
}
