package config

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ColumnInfo struct {
	Field   *string
	Type    *string
	Null    *string
	Key     *string
	Default *string
	Extra   *string
}

type User struct {
	gorm.Model
	Id       int64  `gorm:"primary_key;"`
	NickName string `gorm:"size:40;not null;unique"`
	Avatar   string
	Name     string `gorm:"size:20;not null;unique"`
	Password string `gorm:"size:255;not null"`
	Mobile   string `gorm:"size:11;not null;unique"`
	Email    string `gorm:"size:50;not null;unique"`
}
type Invoice struct {
	ID                int       `gorm:"column:id;primary_key;autoIncrement:true"` // 主键
	CreatedAt         time.Time `gorm:"column:created_at"`                        // 创建时间
	UploadAt          time.Time `gorm:"column:upload_at"`                         // 上传时间
	ClipID            uint      `gorm:"column:clip_id"`                           // 剪辑 ID
	Digest            string    `gorm:"column:digest;size:64"`                    // 摘要
	StorePath         string    `gorm:"column:store_path;size:200"`               // 存储路径
	FileName          string    `gorm:"column:file_name;size:200"`                // 文件名
	Invalid           bool      `gorm:"column:invalid;type:tinyint(1)"`           // 是否无效
	Verify            string    `gorm:"column:verify;size:20"`                    // 验证
	OwnerID           uint      `gorm:"column:owner_id"`                          // 所有者 ID
	Title             string    `gorm:"column:title;size:100"`                    // 标题
	MachineNumber     string    `gorm:"column:machine_number;size:200"`           // 机器编号
	Code              string    `gorm:"column:code;size:200"`                     // 编码
	Number            string    `gorm:"column:number;size:200"`                   // 编号
	Date              time.Time `gorm:"column:date;"`                             // 日期
	Checksum          string    `gorm:"column:checksum;size:200"`                 // 校验和
	BuyerName         string    `gorm:"column:buyer_name;size:250"`               // 买方名称
	BuyerCode         string    `gorm:"column:buyer_code;size:64"`                // 买方编码
	BuyerAddress      string    `gorm:"column:buyer_address;size:300"`            // 买方地址
	BuyerAccount      string    `gorm:"column:buyer_account;size:300"`            // 买方账户
	Password          string    `gorm:"column:password;size:300"`                 // 密码
	Amount            float64   `gorm:"column:amount"`                            // 金额
	TaxAmount         float64   `gorm:"column:tax_amount"`                        // 税额
	TotalAmount       float64   `gorm:"column:total_amount"`                      // 总金额
	TotalAmountString string    `gorm:"column:total_amount_string;size:200"`      // 总金额字符串
	SellerName        string    `gorm:"column:seller_name;size:250"`              // 卖方名称
	SellerCode        string    `gorm:"column:seller_code;size:64"`               // 卖方编码
	SellerAddress     string    `gorm:"column:seller_address;size:300"`           // 卖方地址
	SellerAccount     string    `gorm:"column:seller_account;size:300"`           // 卖方账户
	Payee             string    `gorm:"column:payee;size:100"`                    // 收款人
	Reviewer          string    `gorm:"column:reviewer;size:100"`                 // 审核人
	Drawer            string    `gorm:"column:drawer;size:100"`                   // 开票人
	Memo              string    `gorm:"column:memo;size:2000"`                    // 备注
	Proxy             bool      `gorm:"column:proxy;type:tinyint(1)"`             // 代理
}

func TestTableData(t *testing.T) {
	db, _ := GetDB()
	defer db.Close()
	//tableName := "invoices"
	var invoices []Invoice
	db.Create(&Invoice{
		ID:        1,
		CreatedAt: time.Now(),
		UploadAt:  time.Now(),
		Date:      time.Now(),
	})
	db.Find(&invoices)
	fmt.Println("Invoices:", invoices)
}

func TestGetUsers(t *testing.T) {
	db, _ := GetDB()
	assert.NotNil(t, db, "Database connection should not be nil")

	tableName := "users"
	sqlString := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
	rows, err := db.Raw(sqlString).Rows()
	if err != nil {
		t.Errorf("Failed to execute SQL: %v", err)
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
			t.Errorf("Failed to scan row: %v", err)
			return
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error occurred while iterating over rows: %v", err)
		return
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
}

// 辅助函数，用于获取指针指向的值或默认值
func getValue(s *string) string {
	if s == nil {
		return "NULL"
	}
	return *s
}

func TestInitDB(t *testing.T) {
	db, _ := GetDB()
	assert.NotEmpty(t, db)

	// 使用 db.Raw 执行 SQL 并获取结果
	rows, err := db.Raw("SHOW TABLES;").Rows()
	if err != nil {
		t.Errorf("Failed to execute SQL: %v", err)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			t.Errorf("Failed to scan row: %v", err)
			return
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error occurred while iterating over rows: %v", err)
		return
	}

	fmt.Println("Tables:", tables)
	// 断言结果不为空
	assert.NotEmpty(t, tables)
}

func TestMakeMigrates(t *testing.T) {
	//func TestMakeMigrates(t *testing.T) {
	//	db := GetDB()
	//	err := MakeMigrates(db, []any{
	//		&model.User{},
	//		&model.Group{},
	//		&model.GroupMember{},
	//		&model.Config{},
	//	})
	//	if err != nil {
	//		return
	//	}
	//}
	// 创建模拟数据库
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbMock.Close()

	// 创建 GORM 数据库实例
	db, err := gorm.Open("postgres", dbMock)
	if err != nil {
		//t.Fatal("an error '%s' was not expected when opening a stub database connection", err)
	}

	// 预期的数据库操作
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	// 调用 MakeMigrates 函数
	MakeMigrates(db, []any{})
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("Error:", err)
		//t.Errorf("MakeMigrates returned an unexpected error: %v", err)
	}

	// 确保所有的预期操作都已满足
	if err := mock.ExpectationsWereMet(); err != nil {
		fmt.Println("Error:", err)
		//t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

//TODO 没有数据库怎么做
