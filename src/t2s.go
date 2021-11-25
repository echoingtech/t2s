package src

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	pk                = "PRI"
	unique            = "UNI"
	nullAble          = "YES"
	autoincr          = "auto_increment"
	on_update         = "on update CURRENT_TIMESTAMP"
	current_timestamp = "CURRENT_TIMESTAMP"
)

//map for converting mysql type to golang types
var typeForMysqlToGo = map[string]string{
	"int":                "int64",
	"integer":            "int",
	"tinyint":            "int",
	"smallint":           "int",
	"mediumint":          "int",
	"bigint":             "int64",
	"int unsigned":       "int",
	"integer unsigned":   "int",
	"tinyint unsigned":   "int",
	"smallint unsigned":  "int",
	"mediumint unsigned": "int",
	"bigint unsigned":    "int",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"json":               "string",
	"date":               "time.Time", // time.Time or string
	"datetime":           "time.Time", // time.Time or string
	"timestamp":          "time.Time", // time.Time or string
	"time":               "time.Time", // time.Time or string
	"float":              "float64",
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string",
	"varbinary":          "string",
}

type Table2Struct struct {
	dsn            string
	savePath       string
	db             *sql.DB
	tables         string
	prefix         string
	err            error
	realNameMethod string
	enableJsonTag  bool
	packageName    string
	tagKey         string
}

func NewTable2Struct() *Table2Struct {
	return &Table2Struct{}
}

func (t *Table2Struct) Dsn(d string) *Table2Struct {
	t.dsn = d
	return t
}

func (t *Table2Struct) TagKey(r string) *Table2Struct {
	t.tagKey = r
	return t
}

func (t *Table2Struct) PackageName(r string) *Table2Struct {
	t.packageName = r
	return t
}

func (t *Table2Struct) RealNameMethod(r string) *Table2Struct {
	t.realNameMethod = r
	return t
}

func (t *Table2Struct) SavePath(p string) *Table2Struct {
	t.savePath = p
	return t
}

func (t *Table2Struct) DB(d *sql.DB) *Table2Struct {
	t.db = d
	return t
}

func (t *Table2Struct) Tables(tab string) *Table2Struct {
	t.tables = tab
	return t
}

func (t *Table2Struct) Prefix(p string) *Table2Struct {
	t.prefix = p
	return t
}

func (t *Table2Struct) EnableJsonTag(p bool) *Table2Struct {
	t.enableJsonTag = p
	return t
}

func (t *Table2Struct) Run() error {

	tables := strings.Split(t.tables, ",")

	// 链接mysql, 获取db对象
	t.dialMysql()
	if t.err != nil {
		return t.err
	}

	for _, inputTableName := range tables {

		fmt.Println(fmt.Sprintf("\n\n ************************************************** start convert %s ************************************************** \n\n", inputTableName))

		tableName := strings.Replace(inputTableName, " ", "", -1)
		tableName = strings.Replace(tableName, "\n", "", -1)

		// 获取表和字段的shcema
		tableColumns, err := t.getColumns(tableName)
		if err != nil {
			return err
		}

		// 包名
		var packageName string
		if t.packageName == "" {
			packageName = "package mysql\n\n"
		} else {
			packageName = fmt.Sprintf("package %s\n\n", t.packageName)
		}

		// 组装struct
		var structContent string
		for tableRealName, item := range tableColumns {
			// 去除前缀
			if t.prefix != "" {
				tableRealName = tableRealName[len(t.prefix):]
			}
			tableName := tableRealName

			switch len(tableName) {
			case 0:
			case 1:
				tableName = strings.ToUpper(tableName[0:1])
			default:
				var str string
				tableNames := regexp.MustCompile("[-_]").Split(tableName, -1)
				for _, name := range tableNames {
					str += strings.ToUpper(name[0:1]) + name[1:]
				}
				tableName = str
			}
			//fmt.Println("convert tables", tableName)
			depth := 1
			structContent += "type " + tableName + " struct {\n"
			for _, v := range item {
				//structContent += tab(depth) + v.ColumnName + " " + v.Type + " " + v.Json + "\n"
				// 字段注释
				var clumnComment string
				if v.ColumnComment != "" {
					clumnComment = fmt.Sprintf(" // %s", v.ColumnComment)
				}
				structContent += fmt.Sprintf("%s%s %s %s%s\n",
					tab(depth), v.ColumnName, v.Type, v.Tag, clumnComment)
			}
			structContent += tab(depth-1) + "}\n\n"

			structContent += fmt.Sprintf("func (%s) %s() string {\n",
				tableName, "TableName")
			structContent += fmt.Sprintf("%sreturn \"%s\"\n",
				tab(depth), tableRealName)
			structContent += "}\n\n"

			// 添加 method 获取真实表名
			if t.realNameMethod != "" {
				structContent += fmt.Sprintf("func (*%s) %s() string {\n",
					tableName, t.realNameMethod)
				structContent += fmt.Sprintf("%sreturn \"%s\"\n",
					tab(depth), tableRealName)
				structContent += "}\n\n"
			}
		}

		// 如果有引入 time.Time, 则需要引入 time 包
		var importContent string
		if strings.Contains(structContent, "time.Time") {
			importContent = "import \"time\"\n\n"
		}

		fmt.Println(packageName)
		fmt.Println(importContent)
		fmt.Println(structContent)

		fmt.Println(fmt.Sprintf("\n\n ************************************************** end convert %s ************************************************** \n\n", inputTableName))

		// 写入文件struct
		var savePath = t.savePath
		// 是否指定保存路径
		if savePath == "" {
			continue
		}

		err = func(inputTableName string) error {
			if filepath.Ext(savePath) != ".go" {
				savePath = filepath.Join(savePath, inputTableName+".go")
			}

			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			filePath := filepath.Join(wd, savePath)

			// 确保目录存在
			if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}

			f, err := os.Create(filePath)
			if err != nil {
				fmt.Println("Can not write file")
				return err
			}

			defer func() {
				_ = f.Close()
			}()

			_, err = f.WriteString(packageName + importContent + structContent)
			if err != nil {
				return err
			}

			cmd := exec.Command("gofmt", "-w", filePath)
			if stderr, err := cmd.CombinedOutput(); err != nil {
				fmt.Errorf("gofmt failed: %s", string(stderr))
				return err
			}

			return nil
		}(inputTableName)

		if err != nil {
			fmt.Println("-------------------")
			fmt.Println(err)
			fmt.Println("-------------------")
			panic(err)
		}

	}

	return nil
}

func (t *Table2Struct) dialMysql() {
	if t.db == nil {
		if t.dsn == "" {
			t.err = errors.New("dsn数据库配置缺失")
			return
		}
		t.db, t.err = sql.Open("mysql", t.dsn)
	}
	return
}

type column struct {
	ColumnName    string
	NullAble      string
	Key           string
	Type          string
	Nullable      string
	TableName     string
	ColumnComment string
	Tag           string
	ColumnType    string
	Extra         string
	ColumnDefault *string
}

// Function for fetching schema definition of passed tables
func (t *Table2Struct) getColumns(table string) (tableColumns map[string][]column, err error) {
	tableColumns = make(map[string][]column)
	// sql
	var sqlStr = `
		SELECT 
		COLUMN_NAME, 
		DATA_TYPE, 
		IS_NULLABLE,
		TABLE_NAME,
		COLUMN_COMMENT,
		COLUMN_TYPE,
		COLUMN_KEY,
		IS_NULLABLE,
		EXTRA,
		COLUMN_DEFAULT
		FROM information_schema.COLUMNS 
		WHERE table_schema = DATABASE()`
	// 是否指定了具体的table
	//if t.tables != "" {
	//	sqlStr += fmt.Sprintf(" AND TABLE_NAME = '%s'", t.prefix+t.tables)
	//}
	sqlStr += fmt.Sprintf(" AND TABLE_NAME = '%s'", t.prefix+table)
	// sql排序
	sqlStr += " order by TABLE_NAME asc, ORDINAL_POSITION asc"

	rows, err := t.db.Query(sqlStr)
	if err != nil {
		fmt.Println("Error reading tables information: ", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		col := column{}
		err = rows.Scan(&col.ColumnName, &col.Type, &col.Nullable, &col.TableName, &col.ColumnComment, &col.ColumnType, &col.Key, &col.NullAble, &col.Extra, &col.ColumnDefault)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//fmt.Println("=====>", col)
		//col.Json = strings.ToLower(col.ColumnName)
		col.Tag = col.ColumnName
		col.ColumnComment = col.ColumnComment
		col.ColumnName = t.camelCase(col.ColumnName)
		col.Type = typeForMysqlToGo[col.Type]
		col.Tag = strings.ToLower(col.Tag)

		//if t.enableJsonTag {
		//	col.Tag = fmt.Sprintf("`%s:\"%s\" json:\"%s\"`", "xorm", col.Tag, col.Tag)
		//} else {
		//	col.Tag = fmt.Sprintf("`%s:\"%s %s '%s'\"`", "xorm", col.ColumnType, col.Tag)
		//}
		colStr := fmt.Sprintf("`%s:\"%s", "xorm", col.ColumnType)
		if col.Key == pk {
			colStr = colStr + " pk "
		}
		if col.Key == unique {
			colStr = colStr + " unique "
		}
		if col.Extra == autoincr {
			colStr = colStr + " autoincr "
		}
		if col.NullAble == nullAble {
			colStr = colStr + " null "
		} else {
			colStr = colStr + " notnull "
		}

		if col.Extra == on_update {
			colStr = colStr + " updated "
		}

		if col.ColumnDefault != nil && *col.ColumnDefault == current_timestamp && col.Extra == "" {
			colStr = colStr + " created "
		}

		colStr = colStr + fmt.Sprintf("'%s'\"`", col.Tag)

		col.Tag = colStr
		//col.Tag = fmt.Sprintf(" '%s'\"`", "xorm", col.ColumnType, col.Tag)

		if _, ok := tableColumns[col.TableName]; !ok {
			tableColumns[col.TableName] = []column{}
		}
		tableColumns[col.TableName] = append(tableColumns[col.TableName], col)
	}
	return
}

func (t *Table2Struct) camelCase(str string) string {
	// 是否有表前缀, 设置了就先去除表前缀
	if t.prefix != "" {
		str = strings.Replace(str, t.prefix, "", 1)
	}
	var text string
	for _, p := range strings.Split(str, "_") {
		// 字段首字母大写的同时, 是否要把其他字母转换为小写
		switch len(p) {
		case 0:
		case 1:
			text += strings.ToUpper(p[0:1])
		default:
			text += strings.ToUpper(p[0:1]) + p[1:]
		}
	}
	return text
}
func tab(depth int) string {
	return strings.Repeat("\t", depth)
}
